package main

import (
	"log"
	"os"
	"github.com/gin-gonic/gin"
	api "github.com/n4vxn/Hotel-Hop/api/handlers"
	"github.com/n4vxn/Hotel-Hop/api/middleware"
	"github.com/n4vxn/Hotel-Hop/db"
	"github.com/n4vxn/Hotel-Hop/scripts"
)

func main() {
	dbConn, err := db.NewDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	var (
		userStore    = db.NewPostgresUserStore(dbConn)
		hotelStore   = db.NewPostgresHotelStore(dbConn)
		roomStore    = db.NewPostgresRoomStore(dbConn)
		bookingStore = db.NewPostgresBookingStore(dbConn)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		authHandler    = api.NewAuthHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
	)

	if err := userStore.InitUsersTable(); err != nil {
		log.Fatalf("Failed to initialize users store: %v", err)
	}
	if err := hotelStore.InitHotelsTable(); err != nil {
		log.Fatalf("Failed to initialize hotels store: %v", err)
	}
	if err := roomStore.InitRoomsTable(); err != nil {
		log.Fatalf("Failed to initialize rooms store: %v", err)
	}
	if err := bookingStore.InitBookingsTable(); err != nil {
		log.Fatalf("Failed to initialize bookings store: %v", err)
	}
	//seeding
	if err := scripts.SeedHotelStore(dbConn); err != nil {
		log.Fatalf("error seeding hotel store: %v", err)
	}
	if err := scripts.SeedRoomStore(dbConn); err != nil {
		log.Fatalf("error seeding room store: %v", err)
	}

	var (
		router = gin.Default()
		apiv1  = router.Group("/api/v1/", middleware.JWTAuthMiddleware(userStore))
		auth   = router.Group("/api")
		admin  = apiv1.Group("/admin", middleware.AdminAuth)
	)

	//auth
	auth.POST("/auth", authHandler.HandleAuthenticate)
	auth.POST("/users", userHandler.HandlePostUsers)

	//Versioned API routes
	//user handler
	apiv1.GET("/user/:id", userHandler.HandlerGetUser)
	apiv1.GET("/users", userHandler.HandlerGetUsers)
	apiv1.DELETE("/user/:id", userHandler.HandleDeleteUsers)
	apiv1.PUT("/user/:id", userHandler.HandleUpdateUsers)

	//hotel handlers
	apiv1.POST("/hotel", hotelHandler.HandleCreateHotel)
	apiv1.GET("/hotel", hotelHandler.HandleGetHotels)
	apiv1.GET("/hotel/:id", hotelHandler.HandleGetRooms)
	apiv1.GET("/hotel/:id/rooms", hotelHandler.HandleGetHotelByID)

	//room handler
	apiv1.POST("/hotel/:id/book", roomHandler.HandleBookRoom)

	//booking handlers
	apiv1.GET("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.POST("/booking/:id/cancel", bookingHandler.HandleCancelBookings)

	//admin
	admin.GET("/booking", bookingHandler.HandleGetBookings)
	admin.POST("/booking/:id/cancel", bookingHandler.HandleCancelBookings)

	listenAddr := os.Getenv("LISTEN_ADDR")
	router.Run(listenAddr)
}
