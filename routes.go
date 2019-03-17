package main

import "net/http"

func InitRoutes() {
	// WrappedLoginHandler := getWrappedLoginHandler()
	WrappedIOTDataHandler := getWrappedIOTDataHandler()

	// http.HandleFunc("/api/login", WrappedLoginHandler)
	// http.HandleFunc("/api/login", loggingMiddleware(ExampleLoginHandler))
	// http.HandleFunc("/api/logout", applyMiddlewares(LogoutHandler))
	// http.HandleFunc("/api/authcheck", applyMiddlewares(AuthcheckHandler))

	http.HandleFunc("/api/dashboard", applyMiddlewares(DashboardHandler))

	http.HandleFunc("/api/upgrade", applyMiddlewares(UpgradeHandler))
	http.HandleFunc("/api/upgradecheck", applyMiddlewares(UpgradeCheckHandler))

	http.HandleFunc("/metrics", WrappedIOTDataHandler)
}
