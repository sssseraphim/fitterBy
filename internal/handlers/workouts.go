package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sssseraphim/fitterBy/internal/database"
)

type WorkoutHandler struct {
	DB *database.Queries
}

type Workout struct {
	ID           uuid.UUID     `json:"id"`
	UserId       uuid.UUID     `json:"user_id"`
	ProgramDayId uuid.UUID     `json:"program_day_id"`
	CreatedAt    time.Time     `json:"created_at"`
	Lifts        []WorkoutLift `json:"lifts"`
}

type WorkoutLift struct {
	UserId     uuid.UUID `json:"user_id"`
	WorkoutId  uuid.UUID `json:"workout_id"`
	ExerciseId uuid.UUID `json:"exercise_id"`
	Weight     int       `json:"weight"`
	Sets       int       `json:"sets"`
	Reps       int       `json:"reps"`
	LiftOrder  int       `json:"lift_order"`
}

func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromContext(r)
	var req struct {
		ProgramID    uuid.UUID     `json:"program_id"`
		ProgramDayID uuid.UUID     `json:"program_day_id"`
		Lifts        []WorkoutLift `json:"lifts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, "failed to decode", err)
		return
	}
	workout, err := h.DB.CreateWorkout(r.Context(), database.CreateWorkoutParams{
		UserID:       userId,
		ProgramDayID: req.ProgramDayID,
	})
	if err != nil {
		respondWithError(w, 500, "failed to create workout", err)
		return
	}
	for _, l := range req.Lifts {
		_, err := h.DB.CreateWorkoutLift(r.Context(), database.CreateWorkoutLiftParams{
			UserID:     userId,
			WorkoutID:  workout.ID,
			ExerciseID: l.ExerciseId,
			Sets:       int32(l.Sets),
			Reps:       int32(l.Reps),
			LiftOrder:  int32(l.LiftOrder),
			Weight:     int32(l.Weight),
		})
		if err != nil {
			respondWithError(w, 500, "failed to create workout", err)
			return
		}
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"success": "success"})
}

func (h *WorkoutHandler) HandleGetMyWorkouts(w http.ResponseWriter, r *http.Request) {
	userId := userIdFromContext(r)
	workouts, err := h.DB.GetUsersWorkouts(r.Context(), userId)
	if err != nil {
		respondWithError(w, 500, "failed to get workouts", err)
		return
	}
	var resp struct {
		Workouts []Workout `json:"workouts"`
	}
	for _, w := range workouts {
		resp.Workouts = append(resp.Workouts, Workout{
			ID:           w.ID,
			CreatedAt:    w.CreatedAt.Time,
			ProgramDayId: w.ProgramDayID,
		})
	}
	respondWithJSON(w, 200, resp)
}

func (h *WorkoutHandler) HandleGetWorkout(w http.ResponseWriter, r *http.Request) {
	workoutIdString := r.PathValue("workout_id")
	if workoutIdString == "" {
		respondWithError(w, 400, "workout id required", errors.New("no id"))
		return
	}
	workoutId, err := uuid.Parse(workoutIdString)
	if err != nil {
		respondWithError(w, 400, "wrong workout id", err)
		return
	}
	workout, err := h.DB.GetWorkoutByID(r.Context(), workoutId)
	if err != nil {
		respondWithError(w, 500, "failed to get workout", err)
		return
	}
	resp := Workout{
		ID:           workout.ID,
		ProgramDayId: workout.ProgramDayID,
		CreatedAt:    workout.CreatedAt.Time,
	}
	lifts, err := h.DB.GetWorkoutLifts(r.Context(), workout.ID)
	if err != nil {
		respondWithError(w, 500, "failed to get workout", err)
		return
	}
	for _, l := range lifts {
		resp.Lifts = append(resp.Lifts, WorkoutLift{
			ExerciseId: l.ExerciseID,
			Weight:     int(l.Weight),
			Sets:       int(l.Sets),
			Reps:       int(l.Reps),
			LiftOrder:  int(l.LiftOrder),
		})
	}
	respondWithJSON(w, 200, resp)

}
