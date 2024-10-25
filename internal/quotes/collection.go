package quotes

import "math/rand"

type Collection struct {
	quotes []string
}

func NewCollection() *Collection {
	return &Collection{
		quotes: []string{
			"The only limit to our realization of tomorrow is our doubts of today.",
			"Success is not final, failure is not fatal: It is the courage to continue that counts.",
			"It does not matter how slowly you go as long as you do not stop.",
			"In the middle of difficulty lies opportunity.",
			"Do not wait for the perfect moment, take the moment and make it perfect.",
			"What lies behind us and what lies before us are tiny matters compared to what lies within us.",
			"The only way to do great work is to love what you do.",
			"Happiness is not something ready-made. It comes from your own actions.",
			"The best time to plant a tree was 20 years ago. The second best time is now.",
			"Your time is limited, don't waste it living someone else's life.",
			"The harder you work for something, the greater you'll feel when you achieve it.",
			"Don't watch the clock; do what it does. Keep going.",
			"Dream big and dare to fail.",
			"You don't have to be great to start, but you have to start to be great.",
			"What we achieve inwardly will change outer reality.",
			"Difficulties strengthen the mind, as labor does the body.",
			"Believe you can and you're halfway there.",
			"Challenges are what make life interesting and overcoming them is what makes life meaningful.",
			"Opportunities don't happen, you create them.",
		},
	}
}

func (c *Collection) GetRandomQuote() string {
	return c.quotes[rand.Intn(len(c.quotes))]
}
