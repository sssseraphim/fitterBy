package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sssseraphim/fitterBy/internal/database"
)

type ProgramHandler struct {
	DB *database.Queries
}

type Exercise struct {
	ID          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	AuthorName  string    `json:"author_name"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MediaUrls   []string  `json:"media_urls"`
	CreatedAt   time.Time `json:"created_at"`
}

type Lift struct {
	ExerciseId   uuid.UUID `json:"exercise_id"`
	ExerciseName string    `json:"exercise_name"`
	Sets         int       `json:"sets"`
	Reps         int       `json:"reps"`
	Description  string    `json:"description"`
	Order        int       `json:"order"`
}

type Day struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Lifts       []Lift    `json:"lifts"`
	Order       int       `json:"order"`
}

type Program struct {
	ID          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	AuthorName  string    `json:"author_name"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MediaUrls   []string  `json:"media_urls"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Visibility  string    `json:"visibility"`
	Days        []Day     `json:"days"`
}

type UserProgram struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	ProgramID  uuid.UUID `json:"program_id"`
	CreatedAt  time.Time `json:"created_at"`
	CurrentDay int       `json:"current_day"`
	Status     string    `json:"status"`
}

func (h *ProgramHandler) HandleCreateExercise(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromContext(r)
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		MediaUrls   []string `json:"media_urls"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, "failed to decode request", err)
		return
	}
	exercise, err := h.DB.CreateExercise(r.Context(), database.CreateExerciseParams{
		UserID:      userId,
		Name:        req.Name,
		Description: req.Description,
		MediaUrls:   req.MediaUrls,
	})
	if err != nil {
		respondWithJSON(w, 500, "failed to create an exercise")
		return
	}
	respondWithJSON(w, http.StatusCreated, Exercise{
		ID:          exercise.ID,
		UserId:      exercise.UserID,
		Name:        exercise.Name,
		Description: exercise.Description,
		MediaUrls:   exercise.MediaUrls,
		CreatedAt:   exercise.CreatedAt,
	})
}

func (h *ProgramHandler) HandleGetExerciseById(w http.ResponseWriter, r *http.Request) {
	exerciseIdString := r.PathValue("exercise_id")
	if exerciseIdString == "" {
		respondWithError(w, 400, "post id required", errors.New("no id"))
		return
	}
	exerciseId, err := uuid.Parse(exerciseIdString)
	if err != nil {
		respondWithError(w, 400, "wrong program id", err)
		return
	}

	exercise, err := h.DB.GetExerciseById(r.Context(), exerciseId)
	if err != nil {
		respondWithJSON(w, 404, "failed to get exercise")
		return
	}
	respondWithJSON(w, http.StatusOK, Exercise{
		ID:          exercise.ID,
		UserId:      exercise.UserID,
		AuthorName:  exercise.AuthorName.String,
		Name:        exercise.Name,
		Description: exercise.Description,
		MediaUrls:   exercise.MediaUrls,
		CreatedAt:   exercise.CreatedAt,
	})
}

func (h *ProgramHandler) HandleGetExercises(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Exercises []Exercise `json:"exercises"`
	}
	exercises, err := h.DB.GetExercises(r.Context())
	if err != nil {
		respondWithError(w, 500, "failed to get exercises", err)
		return
	}
	for _, exercise := range exercises {
		response.Exercises = append(response.Exercises, Exercise{
			ID:          exercise.ID,
			UserId:      exercise.UserID,
			AuthorName:  exercise.AuthorName.String,
			Name:        exercise.Name,
			Description: exercise.Description,
			MediaUrls:   exercise.MediaUrls,
			CreatedAt:   exercise.CreatedAt,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (h *ProgramHandler) HandleCreateProgram(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromContext(r)
	var req Program
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err), err)
		return
	}
	program, err := h.DB.CreateProgram(r.Context(), database.CreateProgramParams{
		Name:        req.Name,
		UserID:      userId,
		Description: req.Description,
		MediaUrls:   req.MediaUrls,
		Visibility:  req.Visibility,
	})
	if err != nil {
		respondWithError(w, 500, "failed to create program", err)
		return
	}
	for _, day := range req.Days {
		d, err := h.DB.CreateProgramDay(r.Context(), database.CreateProgramDayParams{
			ProgramID:   program.ID,
			Name:        day.Name,
			Description: day.Description,
			DayOrder:    int32(day.Order),
		})
		if err != nil {
			respondWithError(w, 500, "failed to create program", err)
			return
		}
		for _, lift := range day.Lifts {
			_, err := h.DB.CreateProgramLift(r.Context(), database.CreateProgramLiftParams{
				ProgramDayID: d.ID,
				ExerciseID:   lift.ExerciseId,
				Description:  lift.Description,
				LiftOrder:    int32(lift.Order),
				Sets:         int32(lift.Sets),
				Reps:         int32(lift.Reps),
			})
			if err != nil {
				respondWithError(w, 500, "failed to create program", err)
				return
			}
		}
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"success": "success"})
}

func (h *ProgramHandler) HandleGetPrograms(w http.ResponseWriter, r *http.Request) {
	programs, err := h.DB.GetPrograms(r.Context())
	if err != nil {
		respondWithError(w, 500, "failed to get programs", err)
		return
	}
	var resp struct {
		Programs []Program `json:"programs"`
	}
	for _, p := range programs {
		resp.Programs = append(resp.Programs, Program{
			ID:          p.ID,
			UserId:      p.UserID,
			Name:        p.Name,
			AuthorName:  p.AuthorName.String,
			Description: p.Description,
			MediaUrls:   p.MediaUrls,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Visibility:  p.Visibility,
		})
	}
	respondWithJSON(w, 200, resp)
}

func (h *ProgramHandler) HandleGetProgram(w http.ResponseWriter, r *http.Request) {
	programIdString := r.PathValue("program_id")
	if programIdString == "" {
		respondWithError(w, 400, "program id required", errors.New("no id"))
		return
	}
	programId, err := uuid.Parse(programIdString)
	if err != nil {
		respondWithError(w, 400, "wrong program id", err)
		return
	}
	fmt.Println(programId)
	program, err := h.DB.GetProgram(r.Context(), programId)
	if err != nil {
		respondWithError(w, 404, "failed to find program", err)
		return
	}
	resp := Program{
		ID:          program.ID,
		UserId:      program.UserID,
		Name:        program.Name,
		AuthorName:  program.AuthorName.String,
		Description: program.Description,
		MediaUrls:   program.MediaUrls,
		CreatedAt:   program.CreatedAt,
		UpdatedAt:   program.UpdatedAt,
		Visibility:  program.Visibility,
	}
	days, err := h.DB.GetProgramDays(r.Context(), programId)
	if err != nil {
		respondWithError(w, 500, "failed to find program days", err)
		return
	}
	for _, d := range days {
		respDay := Day{
			ID:          d.ID,
			Name:        d.Name,
			Description: d.Description,
			Order:       int(d.DayOrder),
		}
		lifts, err := h.DB.GetProgramDayLifts(r.Context(), d.ID)
		if err != nil {
			respondWithError(w, 500, "failed to find program day lifts", err)
			return
		}
		for _, l := range lifts {
			respDay.Lifts = append(respDay.Lifts, Lift{
				ExerciseId:  l.ExerciseID,
				Description: l.Description,
				Sets:        int(l.Sets),
				Reps:        int(l.Reps),
				Order:       int(l.LiftOrder),
			})
		}
		resp.Days = append(resp.Days, respDay)
	}
	respondWithJSON(w, 200, resp)

}

func (h *ProgramHandler) HandleSubscribeToProgram(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromContext(r)
	programIdString := r.PathValue("program_id")
	if programIdString == "" {
		respondWithError(w, 400, "program id required", errors.New("no id"))
		return
	}
	programId, err := uuid.Parse(programIdString)
	if err != nil {
		respondWithError(w, 400, "wrong program id", err)
		return
	}
	err = h.DB.SubscribeToProgram(r.Context(), database.SubscribeToProgramParams{
		UserID:    userId,
		ProgramID: programId})
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("failed to subscribe: %v", err), err)
		return
	}
	respondWithJSON(w, 200, map[string]string{"success": "success"})
}

func (h *ProgramHandler) HandleGetSubscribedPrograms(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromContext(r)
	programs, err := h.DB.GetUserSubscribedPrograms(r.Context(), userId)
	if err != nil {
		respondWithError(w, 500, "failed to get programs", err)
		return
	}
	var resp struct {
		Programs []UserProgram `json:"programs"`
	}
	for _, p := range programs {
		resp.Programs = append(resp.Programs, UserProgram{
			ID:         p.ID,
			Name:       p.Name.String,
			ProgramID:  p.ProgramID,
			CreatedAt:  p.CreatedAt.Time,
			CurrentDay: int(p.CurrentDayOrder.Int32),
			Status:     p.Status.String,
		})
	}
	respondWithJSON(w, 200, resp)
}
