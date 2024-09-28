package main

import (
	nn_adapter "annotater/internal/bl/NN/NNAdapter"
	nn_model_handler "annotater/internal/bl/NN/NNAdapter/NNmodelhandler"
	annot_service "annotater/internal/bl/annotationService"
	annot_repo_adapter "annotater/internal/bl/annotationService/annotattionRepo/anotattionRepoAdapter"
	annot_type_service "annotater/internal/bl/anotattionTypeService"
	annot_type_repo_adapter "annotater/internal/bl/anotattionTypeService/anottationTypeRepo/anotattionTypeRepoAdapter"
	auth_service "annotater/internal/bl/auth"
	document_service "annotater/internal/bl/documentService"
	doc_data_repo_adapter "annotater/internal/bl/documentService/documentDataRepo/documentDataRepo"
	document_repo_adapter "annotater/internal/bl/documentService/documentMetaDataRepo/documentMetaDataRepoAdapter"
	rep_data_repo_adapter "annotater/internal/bl/documentService/reportDataRepo/reportDataRepoAdapter"
	rep_creator_service "annotater/internal/bl/reportCreatorService"
	report_creator "annotater/internal/bl/reportCreatorService/reportCreator"
	service "annotater/internal/bl/userService"
	user_repo_adapter "annotater/internal/bl/userService/userRepo/userRepoAdapter"
	"annotater/internal/config"
	annot_handler "annotater/internal/http-server/handlers/annot"
	annot_type_handler "annotater/internal/http-server/handlers/annotType"
	auth_handler "annotater/internal/http-server/handlers/auth"
	document_handler "annotater/internal/http-server/handlers/document"
	user_handler "annotater/internal/http-server/handlers/user"
	logger_setup "annotater/internal/logger"
	"annotater/internal/middleware/access_middleware"
	"annotater/internal/middleware/auth_middleware"
	models_da "annotater/internal/models/modelsDA"
	auth_utils "annotater/internal/pkg/authUtils"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// andrew1 2
// admin admin
// control control

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models_da.Document{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models_da.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models_da.MarkupType{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models_da.Markup{})
	if err != nil {
		return err
	}
	return nil
}

func main() {

	config := config.MustLoad()
	postgresConStr := config.Database.GetGormConnectStr()
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: postgresConStr}),
		&gorm.Config{TranslateError: true,
			Logger: logger.Default.LogMode(logger.Silent)})

	log := logger_setup.Setuplog(config)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = migrate(db)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//auth service
	userRepo := user_repo_adapter.NewUserRepositoryAdapter(db)
	hasher := auth_utils.NewPasswordHashCrypto()
	tokenHandler := auth_utils.NewJWTTokenHandler()
	authService := auth_service.NewAuthService(log, userRepo, hasher, tokenHandler, auth_service.SECRET)

	//annot service
	annotRepo := annot_repo_adapter.NewAnotattionRepositoryAdapter(db)
	annotService := annot_service.NewAnnotattionService(log, annotRepo)

	//annotType service
	annotTypeRepo := annot_type_repo_adapter.NewAnotattionTypeRepositoryAdapter(db)
	annotTypeService := annot_type_service.NewAnotattionTypeService(log, annotTypeRepo)

	//document service
	//setting up NN
	modelhandler := nn_model_handler.NewHttpModelHandler(log, config.Model.Route)
	model := nn_adapter.NewDetectionModel(modelhandler)

	reportCreator := report_creator.NewPDFReportCreator(config.ReportCreatorPath)
	reportCreatorService := rep_creator_service.NewDocumentService(log, model, annotTypeRepo, reportCreator)

	documentStorage := doc_data_repo_adapter.NewDocumentRepositoryAdapter(config.DocumentPath, config.DocumentExt)

	reportStorage := rep_data_repo_adapter.NewDocumentRepositoryAdapter(config.ReportPath, config.ReportExt)

	documentRepo := document_repo_adapter.NewDocumentRepositoryAdapter(db)
	documentService := document_service.NewDocumentService(log, documentRepo, documentStorage, reportStorage, reportCreatorService)

	//userService 0_0
	userService := service.NewUserService(log, userRepo)

	//handlers
	userHandler := user_handler.NewDocumentHandler(log, userService)
	documentHandlerV1 := document_handler.NewDocumentHandlerV1(log, documentService)
	documentHandlerV2 := document_handler.NewDocumentHandlerV2(log, documentService)

	annotHandlerV1 := annot_handler.NewAnnotHandlerV1(log, annotService)
	annotHandlerV2 := annot_handler.NewAnnotHandlerV2(log, annotService)

	annotTypeHandler := annot_type_handler.NewAnnotTypehandler(log, annotTypeService)

	authHandlerV1 := auth_handler.NewAuthHandlerV1(log, authService)
	authHandlerV2 := auth_handler.NewAuthHandlerV2(log, authService)

	//auth service
	router := chi.NewRouter()
	//router.Use(middleware.Logger)

	authMiddleware := auth_middleware.NewJwtAuthMiddleware(log, auth_service.SECRET, tokenHandler)
	accesMiddleware := access_middleware.NewAccessMiddleware(log, userService)

	router.Use(middleware.Logger)

	router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Group(func(r chi.Router) { // group for which auth middleware is required
				r.Use(authMiddleware.MiddlewareFunc)

				// Document
				r.Route("/document", func(r chi.Router) {
					r.Post("/report", documentHandlerV1.CreateReport())
					r.Get("/getDocument", documentHandlerV1.GetDocumentByID())
					r.Get("/getReport", documentHandlerV1.GetReportByID())
					r.Get("/getDocumentsMeta", documentHandlerV1.GetDocumentsMetaData())
				})

				// AnnotType
				r.Route("/annotType", func(r chi.Router) {
					r.Use(accesMiddleware.ControllersAndHigherMiddleware) // apply the desired middleware here

					adminOnlyAnnotTypes := r.Group(nil)
					adminOnlyAnnotTypes.Use(accesMiddleware.AdminOnlyMiddleware)

					r.Post("/add", annotTypeHandler.AddAnnotType())
					r.Get("/get", annotTypeHandler.GetAnnotType())

					r.Get("/creatorID", annotTypeHandler.GetAnnotTypesByCreatorID())

					r.Get("/gets", annotTypeHandler.GetAnnotTypesByIDs())

					adminOnlyAnnotTypes.Delete("/delete", annotTypeHandler.DeleteAnnotType())
					r.Get("/getsAll", annotTypeHandler.GetAllAnnotTypes())

				})
				//Annot
				r.Route("/annot", func(r chi.Router) {
					r.Use(accesMiddleware.ControllersAndHigherMiddleware)
					//adminOnlyAnnots := r.Group(nil)
					//adminOnlyAnnots.Use(accesMiddleware.AdminOnlyMiddleware)

					r.Post("/add", annotHandlerV1.AddAnnot())
					r.Get("/get", annotHandlerV1.GetAnnot())
					r.Get("/creatorID", annotHandlerV1.GetAnnotsByUserID())

					r.Delete("/delete", annotHandlerV1.DeleteAnnot())
					r.Get("/getsAll", annotHandlerV1.GetAllAnnots())
				})
				//user
				r.Route("/user", func(r chi.Router) {
					r.Use(accesMiddleware.AdminOnlyMiddleware)
					r.Post("/role", userHandler.ChangeUserPerms())
					r.Get("/getUsers", userHandler.GetAllUsers())
				})

			})

			//auth, no middleware is required
			router.Post("/user/SignUp", authHandlerV1.SignUp())
			router.Post("/user/SignIn", authHandlerV1.SignIn())
		})

		r.Route("/v2", func(r chi.Router) {
			r.Group(func(r chi.Router) { // group for which auth middleware is required
				r.Use(authMiddleware.MiddlewareFunc)

				// Document
				r.Post("/documents", documentHandlerV2.CreateReport())
				r.Get("/documents", documentHandlerV2.GetDocumentsMetaData())
				r.Get("/documents/{id}", documentHandlerV2.GetDocumentByID())

				// Reports
				r.Get("/documents/{id}/reports", documentHandlerV2.GetReportByID())

				// AnnotTypes
				r.With(accesMiddleware.ControllersAndHigherMiddleware).Post("/anottationTypes", nil) //smth
				r.With(accesMiddleware.ControllersAndHigherMiddleware).Get("/anottationTypes", nil)
				r.With(accesMiddleware.ControllersAndHigherMiddleware).Get("/anottationTypes/{id}", nil)

				// Annots
				r.With(accesMiddleware.ControllersAndHigherMiddleware).Post("/anottations", annotHandlerV2.AddAnnot()) //smth
				r.With(accesMiddleware.ControllersAndHigherMiddleware).Get("/anottations", annotHandlerV2.GetAllAnnots())
				r.With(accesMiddleware.ControllersAndHigherMiddleware).Get("/anottations/{id}", annotHandlerV2.GetAnnot())
				r.With(accesMiddleware.ControllersAndHigherMiddleware).Delete("/anottations/{id}", annotHandlerV2.DeleteAnnot())

				// Users
				r.With(accesMiddleware.AdminOnlyMiddleware).Patch("/users/{id}", nil)
				r.With(accesMiddleware.AdminOnlyMiddleware).Get("/users", nil)
			})

			//auth, no middleware is required
			r.Post("/auth", authHandlerV2.Auth())
			r.Post("/register", authHandlerV2.Register())
		})
	})

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         config.Addr,
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("error with server")
		}
	}()

	<-done
}
