package http

func (s *Server) initDiseasesRoutes() {
	diseases := s.server.Group("/api/v1/diseases", s.identifyUser)
	// diseases.GET("/info")
}
