package v1

type GetCollegeRequest struct {
	CollegeId uint `json:"college_id"`
}

type CollegeResponseData struct {
	CollegeId   uint   `json:"college_id"`
	CollegeName string `json:"college_name"`
	Description string `json:"description"`
}

type GetCollegeListDataResponse struct {
	CollegeList []*CollegeResponseData `json:"college_list"`
}
