package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

const (
	Excellent  = "Excellent"
	Proficient = "Proficient"
	Decent     = "Decent"
	Deficient  = "Deficient"
	Poor       = "Poor"
)

// ScoreConvertLetter converts category string to letter grade (A-F)
func ScoreConvertLetter(s string) string {
	grades := map[string]string{
		Excellent:  "A",
		Proficient: "B",
		Decent:     "C",
		Deficient:  "C",
		Poor:       "D",
	}
	if grade, ok := grades[s]; ok {
		return grade
	}
	return "F"
}

// ScoreConvertFloat converts category string to numeric score
func ScoreConvertFloat(s string) float64 {
	switch s {
	case Excellent:
		return 95.0
	case Proficient:
		return 85.0
	case Decent:
		return 75.0
	case Deficient:
		return 65.0
	case Poor:
		return 50.0
	default:
		return 0.0
	}
}

// calculateOverallScore takes the four category strings and returns the average numeric score
func calculateOverallScore(engagement, verbal, classwork, wonder string) float64 {
	return (ScoreConvertFloat(engagement) +
		ScoreConvertFloat(verbal) +
		ScoreConvertFloat(classwork) +
		ScoreConvertFloat(wonder)) / 4.0
}

type WeeklyParticipation struct {
	Engagement               string
	VerbalParticipation      string
	ClassworkAndOrganization string
	WonderAndDepthOfInquiry  string
}

func main() {
	var wp WeeklyParticipation

	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions(Excellent, Proficient, Decent, Deficient, Poor)...).
				Title("Engagement").
				Description("On what level were they engaged?").
				Value(&wp.Engagement),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions(Excellent, Proficient, Decent, Deficient, Poor)...).
				Title("Verbal Participation").
				Description("How well did they verbally participate?").
				Value(&wp.VerbalParticipation),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions(Excellent, Proficient, Decent, Deficient, Poor)...).
				Title("Classwork and Organization").
				Description("Were they organized? Did they do what they were told when they were told to do it?").
				Value(&wp.ClassworkAndOrganization),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions(Excellent, Proficient, Decent, Deficient, Poor)...).
				Title("Wonder and Depth of Inquiry").
				Description("Did they posit creative and meaningful questions? Did they share their curiosity meaningfully during the class?").
				Value(&wp.WonderAndDepthOfInquiry),
		),
	).WithAccessible(accessible)

	err := form.Run()
	if err != nil {
		fmt.Println("Uh oh! :", err)
		os.Exit(1)
	}

	// Spinner for fun
	prepareBurger := func() {
		time.Sleep(2 * time.Second)
	}
	_ = spinner.New().Title("Calculating participation score...").Accessible(accessible).Action(prepareBurger).Run()

	// Calculate overall numeric score
	overallScore := calculateOverallScore(
		wp.Engagement,
		wp.VerbalParticipation,
		wp.ClassworkAndOrganization,
		wp.WonderAndDepthOfInquiry,
	)

	var sb strings.Builder

	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}

	verdictText := fmt.Sprintf(
		"VERDICT\n\n"+
			"Engagement: %s\n"+
			"Verbal Participation: %s\n"+
			"Classwork and Organization: %s\n"+
			"Wonder and Depth of Inquiry: %s\n\n"+
			"Overall Participation Score: %s%.1f%%",
		keyword(ScoreConvertLetter(wp.Engagement)),
		keyword(ScoreConvertLetter(wp.VerbalParticipation)),
		keyword(ScoreConvertLetter(wp.ClassworkAndOrganization)),
		keyword(ScoreConvertLetter(wp.WonderAndDepthOfInquiry)),
		keyword(""),
		overallScore,
	)

	fmt.Fprint(&sb,
		"Participation!\n\n",
		lipgloss.NewStyle().Bold(true).Render(verdictText),
	)

	fmt.Println(
		lipgloss.NewStyle().
			Width(60).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Render(sb.String()),
	)
}
