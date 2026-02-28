package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/config"
	"github.com/hacker4257/pet_charity/internal/handler"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/internal/ws"
	"github.com/hacker4257/pet_charity/pkg/event"
)

func Setup(mode string) *gin.Engine {
	hub := ws.NewHub()
	go hub.Run()
	ws.StartNotifySubscriber(hub)

	gin.SetMode(mode)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.AccessLog())
	r.Use(gin.Recovery())

	// 请求体大小限制
	r.MaxMultipartMemory = 8 << 20 // 8MB
	//静态文件服务
	r.Static("/uploads", "./uploads")

	// 统一创建 Repository（每种只创建一个实例）
	userRepo := repository.NewUserRepo()
	orgRepo := repository.NewOrgRepo()
	petRepo := repository.NewPetRepo()
	adoptionRepo := repository.NewAdoptionRepo()
	donationRepo := repository.NewDonationRepo()
	rescueRepo := repository.NewRescueRepo()
	smsRepo := repository.NewSmsRepo()
	tokenRepo := repository.NewTokenRepo()

	txManager := repository.NewTransactionManager()
	// 创建 Service
	userService := service.NewUserService(userRepo, tokenRepo)
	smsService := service.NewSmsService(smsRepo)
	orgService := service.NewOrgService(orgRepo, userRepo)
	petService := service.NewPetService(petRepo, orgRepo)
	adoptionService := service.NewAdoptionService(adoptionRepo, petRepo, orgRepo, txManager)
	donationService := service.NewDonationService(donationRepo, orgRepo, petRepo, txManager)
	donationService.StartExpireWorker()

	rescueService := service.NewRescueService(rescueRepo, orgRepo, petRepo)

	notifyRepo := repository.NewNotificationRepo()
	notifyService := service.NewNotificationService(notifyRepo)
	notifyHandler := handler.NewNotificationHandler(notifyService)

	// 创建 Handler
	authHandler := handler.NewAuthHandler(userService, smsService, tokenRepo)
	userHandler := handler.NewUserHandler(userService)
	orgHandler := handler.NewOrgHandler(orgService)
	petHandler := handler.NewPetHandler(petService)
	adoptionHandler := handler.NewAdoptionHandler(adoptionService)
	donationHandler := handler.NewDonationHandler(donationService)
	rescueHandler := handler.NewRescueHandler(rescueService)
	adminHandler := handler.NewAdminHandler(orgService, userService, petService, donationService)

	//message
	msgRepo := repository.NewMessageRepo()
	msgService := service.NewMessageService(msgRepo, userRepo)
	chatHandler := handler.NewChatHandler(msgService, hub)

	//注册事件钩子 + 启动事件消费
	activityRepo := repository.NewActivityRepo()
	activityService := service.NewActivityService(activityRepo)
	activityHandler := handler.NewActivityHandler(activityService)

	kafkaCfg := config.Global.Kafka
	event.Init(kafkaCfg.Brokers, kafkaCfg.Topic)

	activityService.RegisterHook()
	// 注册通知事件钩子
	notifyService.RegisterHook()

	cacheRepo := repository.NewCacheRepo()
	event.Subscribe(func(e event.Event) {
		if e.Action == "profile_update" {
			cacheRepo.DeleteUser(e.UserID)
		}
	})

	feedRepo := repository.NewFeedRepo()
	event.Subscribe(func(e event.Event) {
		feedRepo.Push(e.Action, e.UserID)
	})
	event.StartConsumer(kafkaCfg.Brokers, kafkaCfg.Topic, kafkaCfg.GroupID)

	//feed
	feedHandler := handler.NewFeedHandler(feedRepo)

	//favorite

	favoriteRepo := repository.NewFavoriteRepo()
	favoriteService := service.NewFavoritService(favoriteRepo, petRepo)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	//公共api组
	v1 := r.Group("/api/v1")

	//公共读接口限流
	publicRL := middleware.RateLimit("public", 60, time.Minute)

	auth := v1.Group("/auth")
	auth.Use(middleware.RateLimit("auth", 10, time.Minute))
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/sms/send", authHandler.SendSmsCode)
		auth.POST("/sms/login", authHandler.SmsLogin)
	}
	orgsPublic := v1.Group("/organizations")
	orgsPublic.Use(publicRL)
	{
		orgsPublic.GET("", orgHandler.List)
		orgsPublic.GET("/:id", orgHandler.GetByID)
		orgsPublic.GET("/nearby", orgHandler.Nearby)
	}
	petsPublic := v1.Group("/pets")
	petsPublic.Use(publicRL)
	{
		petsPublic.GET("", petHandler.List)
		petsPublic.GET("/:id", petHandler.GetByID)
		petsPublic.GET("/stats", petHandler.PublicStats)
	}

	rescuesPublic := v1.Group("/rescues")
	rescuesPublic.Use(publicRL)
	{
		rescuesPublic.GET("", rescueHandler.List)
		rescuesPublic.GET("/map", rescueHandler.MapData)
		rescuesPublic.GET("/:id", rescueHandler.GetByID)
	}
	//捐赠公开
	v1.GET("/donations/public", donationHandler.ListPublic)
	//排行榜
	v1.GET("/leaderboard", activityHandler.Leaderboard)
	//支付回调（不限流，由支付平台发起）
	v1.POST("/donations/notify/wechat", donationHandler.WechatNotify)
	v1.POST("/donations/notify/alipay", donationHandler.AlipayNotify)

	//状态
	v1.GET("/feed", feedHandler.List)

	//需要登录
	authorized := v1.Group("")
	authorized.Use(middleware.JWTAuth(tokenRepo.GetVersion))
	{
		users := authorized.Group("/users")
		{
			users.GET("/me", authHandler.GetMe)
			users.PUT("/me", userHandler.UpdateProfile)
			users.PUT("/me/password", userHandler.ChangePassword)
			users.PUT("/me/avatar", userHandler.Updateavatar)
		}

		orgs := authorized.Group("/organizations")
		{
			orgs.POST("", orgHandler.Create)
			orgs.PUT("", orgHandler.Update)
			orgs.GET("/mine", orgHandler.GetMine)
		}

		pets := authorized.Group("/pets")
		{
			pets.POST("", petHandler.Create)
			pets.PUT("/:id", petHandler.Update)
			pets.DELETE("/:id", petHandler.Delete)
			pets.POST("/:id/images", petHandler.UploadImage)
			pets.DELETE("/:id/images/:imageId", petHandler.DeleteImage)
			pets.POST("/:id/favorite", favoriteHandler.Toggle)
			pets.GET("/:id/favorite", favoriteHandler.GetStatus)
		}
		//领养
		adoptions := authorized.Group("/adoptions")
		{
			adoptions.POST("", adoptionHandler.Create)
			adoptions.GET("/mine", adoptionHandler.ListMine)
			adoptions.GET("/org", adoptionHandler.ListByOrg)
			adoptions.GET("/:id", adoptionHandler.GetByID)
			adoptions.PUT("/:id/review", adoptionHandler.Review)
			adoptions.PUT("/:id/complete", adoptionHandler.Complete)
		}
		//捐赠
		donations := authorized.Group("/donations")
		{
			donations.POST("", donationHandler.Create)
			donations.GET("/mine", donationHandler.ListMine)
			donations.GET("/:id/status", donationHandler.GetStatus)
		}
		//救助
		rescues := authorized.Group("/rescues")
		{
			rescues.POST("", rescueHandler.Create)
			rescues.PUT("/:id", rescueHandler.Update)
			rescues.POST("/:id/images", rescueHandler.UploadImage)
			rescues.POST("/:id/follow", rescueHandler.AddFollow)
			rescues.POST("/:id/claim", rescueHandler.Claim)
			rescues.PUT("/:id/claim/status", rescueHandler.UpdateClaim)
			rescues.POST("/:id/convert", rescueHandler.ConvertToPet)
		}

		//chat
		chat := authorized.Group("/chat")
		{
			chat.GET("/ws", chatHandler.HandleWS)
			chat.GET("/conversations", chatHandler.Conversations)
			chat.GET("/private/:userId", chatHandler.PrivateHistory)
			chat.GET("/room/:roomId", chatHandler.RoomHistory)
		}

		//通知
		notifications := authorized.Group("/notifications")
		{
			notifications.GET("", notifyHandler.List)
			notifications.GET("/unread", notifyHandler.UnreadCount)
			notifications.PUT("/:id/read", notifyHandler.MarkRead)
			notifications.PUT("/read-all", notifyHandler.MarkAllRead)
		}

		authorized.GET("/leaderboard/me", activityHandler.MyRank)
		authorized.GET("/favorites", favoriteHandler.ListMine)

	}

	//管理后台
	// 管理后台
	admin := v1.Group("/admin")
	admin.Use(middleware.JWTAuth(tokenRepo.GetVersion), middleware.RequiredRole("admin"))
	{
		admin.GET("/stats", adminHandler.Stats)
		admin.GET("/organizations/pending", adminHandler.ListPendingOrgs)
		admin.PUT("/organizations/:id/review", adminHandler.ReviewOrg)
		admin.GET("/donations/stats", donationHandler.Stats)
		admin.GET("/users", adminHandler.ListUsers)
		admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
		admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
	}

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	//warmup geo
	// geo.WarmupGeo()

	return r
}
