package product_review

type Review struct {
	ReviewerName  string `json:"reviewer_name"`
	WrittenReview string `json:"written_review"`
	Rating        int    `json:"rating"`
}
