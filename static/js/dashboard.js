// Dashboard functionality
document.addEventListener('DOMContentLoaded', function() {
    // Load user data and dashboard content
    loadUserData();
    setupNavigation();
    loadDashboardData();
});

// Load user data and update UI
function loadUserData() {
    const user = AuthManager.getUser();
    
    if (user) {
        // Update UI with user data
        document.getElementById('user-name').textContent = user.name;
        document.getElementById('user-greeting').textContent = `Hi, ${user.name.split(' ')[0]}!`;
        document.getElementById('user-type').textContent = user.user_type;
        document.getElementById('welcome-title').textContent = `Welcome, ${user.name.split(' ')[0]}!`;
        
        // Update welcome message based on user type
        const welcomeMessage = user.user_type === 'trainer' 
            ? 'Ready to transform some lives today?' 
            : 'Ready to crush your fitness goals today?';
        document.getElementById('welcome-message').textContent = welcomeMessage;
        
        // Update avatar based on user type
        const avatar = document.getElementById('user-avatar');
        avatar.textContent = user.user_type === 'trainer' ? '💪' : '🚴';
        
    } else {
        // No user data, redirect to login
        AuthManager.redirectToLogin();
    }
}

// Load dashboard data from API
async function loadDashboardData() {
    try {
        const response = await AuthManager.apiCall('/api/dashboard');
        if (response.ok) {
            const data = await response.json();
            updateDashboardUI(data);
        } else {
            console.error('Failed to load dashboard data');
        }
    } catch (error) {
        console.error('Error loading dashboard:', error);
        // Load mock data for demo
        loadMockDashboardData();
    }
}

// Update UI with real data from API
function updateDashboardUI(data) {
    // Update stats
    if (data.stats) {
        document.getElementById('sessions-count').textContent = data.stats.upcoming_sessions || 0;
        document.getElementById('workouts-count').textContent = data.stats.workouts_completed || 0;
        document.getElementById('goals-count').textContent = data.stats.active_goals || 0;
        document.getElementById('streak-count').textContent = data.stats.current_streak || 0;
    }
    
    // Update streak
    if (data.streak) {
        renderStreakBoard(data.streak);
    }
    
    // Update recent workouts
    if (data.recent_workouts) {
        renderRecentWorkouts(data.recent_workouts);
    }
}

// Mock data for demo (remove when API is ready)
function loadMockDashboardData() {
    const user = AuthManager.getUser();
    
    // Mock stats
    document.getElementById('sessions-count').textContent = '2';
    document.getElementById('workouts-count').textContent = '15';
    document.getElementById('goals-count').textContent = '3';
    document.getElementById('streak-count').textContent = '7';
    document.getElementById('current-streak').textContent = '7 days';
    
    // Mock streak data (last 30 days)
    const mockStreak = generateMockStreakData();
    renderStreakBoard(mockStreak);
    
    // Mock recent workouts
    const mockWorkouts = generateMockWorkouts();
    renderRecentWorkouts(mockWorkouts);
}

// Generate mock streak data (like GitHub contributions)
function generateMockStreakData() {
    const streak = [];
    const today = new Date();
    
    for (let i = 29; i >= 0; i--) {
        const date = new Date(today);
        date.setDate(date.getDate() - i);
        
        // Random activity levels for demo
        const activityLevel = Math.floor(Math.random() * 5); // 0-4
        streak.push({
            date: date.toISOString().split('T')[0],
            count: activityLevel,
            level: ['none', 'low', 'medium', 'high', 'max'][activityLevel]
        });
    }
    
    return streak;
}

// Generate mock workouts
function generateMockWorkouts() {
    return [
        {
            id: 1,
            name: "Upper Body Strength",
            type: "Strength Training",
            duration: 45,
            exercises: 8,
            date: new Date().toISOString(),
            calories: 320
        },
        {
            id: 2, 
            name: "Morning Cardio Blast",
            type: "Cardio",
            duration: 30,
            exercises: 5,
            date: new Date(Date.now() - 86400000).toISOString(), // yesterday
            calories: 280
        },
        {
            id: 3,
            name: "Leg Day",
            type: "Strength Training", 
            duration: 60,
            exercises: 6,
            date: new Date(Date.now() - 172800000).toISOString(), // 2 days ago
            calories: 380
        }
    ];
}

// Render streak board (GitHub-style)
function renderStreakBoard(streakData) {
    const streakGrid = document.getElementById('streak-grid');
    streakGrid.innerHTML = '';
    
    streakData.forEach(day => {
        const dayElement = document.createElement('div');
        dayElement.className = `streak-day ${day.level}`;
        dayElement.title = `${day.date}: ${day.count} workouts`;
        streakGrid.appendChild(dayElement);
    });
}

// Render recent workouts list
function renderRecentWorkouts(workouts) {
    const workoutsContainer = document.getElementById('recent-workouts');
    
    if (workouts.length === 0) {
        workoutsContainer.innerHTML = `
            <div class="empty-state">
                <p>No workouts logged yet</p>
                <button class="btn btn-small" onclick="logNewWorkout()">Log Your First Workout</button>
            </div>
        `;
        return;
    }
    
    workoutsContainer.innerHTML = workouts.map(workout => `
        <div class="workout-item">
            <div class="workout-info">
                <h4>${workout.name}</h4>
                <p>${workout.type} • ${workout.exercises} exercises</p>
            </div>
            <div class="workout-stats">
                <div class="workout-duration">${workout.duration} min</div>
                <div class="workout-date">${formatWorkoutDate(workout.date)}</div>
            </div>
        </div>
    `).join('');
}

// Format workout date
function formatWorkoutDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;
    if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
    return date.toLocaleDateString();
}

// Navigation setup (same as before)
function setupNavigation() {
    const navItems = document.querySelectorAll('.nav-item');
    
    navItems.forEach(item => {
        item.addEventListener('click', function(e) {
            e.preventDefault();
            
            // Remove active class from all items
            navItems.forEach(nav => nav.classList.remove('active'));
            // Add active class to clicked item
            this.classList.add('active');
            
            // Hide all sections
            const sections = document.querySelectorAll('.dashboard-section');
            sections.forEach(section => section.classList.remove('active'));
            
            // Show target section
            const targetSection = this.getAttribute('data-section');
            document.getElementById(targetSection).classList.add('active');
        });
    });
}

// Action functions
function bookSession() {
    window.location.href = '/trainers';
}

function logWorkout() {
    logNewWorkout();
}

function findTrainers() {
    window.location.href = '/trainers';
}

function setGoals() {
    showNotification('Goal setting coming soon!', 'info');
}

function logNewWorkout() {
    showNotification('Workout logging feature coming soon!', 'info');
}

// Logout functionality
document.getElementById('logout-btn')?.addEventListener('click', function(e) {
    e.preventDefault();
    
    // Call logout API
    fetch('/api/auth/logout', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${AuthManager.getAccessToken()}`
        }
    }).finally(() => {
        AuthManager.clearTokens();
        window.location.href = '/';
    });
});

