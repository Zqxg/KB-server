package v1

type GetCollegeRequest struct {
	CollegeId int64 `json:"collegeId"`
}

type CollegeResponseData struct {
	CollegeId   uint   `json:"collegeId"`
	CollegeName string `json:"collegeName"`
	Description string `json:"description"`
}
