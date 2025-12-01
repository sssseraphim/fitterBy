// === TOKEN MANAGEMENT SYSTEM ===
class AuthManager {
    static storeTokens(accessToken, refreshToken, user) {
        localStorage.setItem('access_token', accessToken);
        localStorage.setItem('refresh_token', refreshToken);
        localStorage.setItem('user', JSON.stringify(user));
        this.scheduleTokenRefresh();
    }

    static getAccessToken() {
        return localStorage.getItem('access_token');
    }

    static getRefreshToken() {
        return localStorage.getItem('refresh_token');
    }

    static getUser() {
		try {
            const user = localStorage.getItem('user');
            // Check if it's the string "undefined" or actual undefined
            if (!user || user === 'undefined' || user === 'null') {
                return null;
            }
            return JSON.parse(user);
        } catch (error) {
            console.error('Error parsing user data:', error);
            this.clearTokens(); // Clear corrupted data
            return null;
        }
    }

    static clearTokens() {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
    }

    static isAuthenticated() {
        return !!this.getAccessToken();
    }

    static scheduleTokenRefresh() {
        // Refresh token 1 minute before expiry (14 minutes)
        setTimeout(() => {
            this.refreshTokens();
        }, 14 * 60 * 1000);
    }

    static async refreshTokens() {
        const refreshToken = this.getRefreshToken();
        if (!refreshToken) {
            this.redirectToLogin();
            return;
        }

        try {
            const response = await fetch('/api/auth/refresh', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: refreshToken })
            });

            if (response.ok) {
                const data = await response.json();
                this.storeTokens(data.access_token, data.refresh_token, this.getUser());
                console.log('Tokens refreshed automatically');
            } else {
                this.redirectToLogin();
            }
        } catch (error) {
            console.error('Token refresh failed:', error);
            this.redirectToLogin();
        }
    }

    static redirectToLogin() {
        this.clearTokens();
        window.location.href = '/auth';
    }

    // Enhanced fetch with auto token refresh
    static async apiCall(url, options = {}) {
        const token = this.getAccessToken();
        if (token) {
            options.headers = {
                ...options.headers,
                'Authorization': `Bearer ${token}`
            };
        }

        let response = await fetch(url, options);
        
        // If token expired, try to refresh and retry
        if (response.status === 401) {
            await this.refreshTokens();
            const newToken = this.getAccessToken();
            if (newToken) {
                options.headers['Authorization'] = `Bearer ${newToken}`;
                response = await fetch(url, options);
            } else {
                this.redirectToLogin();
                throw new Error('Authentication required');
            }
        }
        
        return response;
    }
}

// === NOTIFICATION SYSTEM ===
function showNotification(message, type = 'info') {
    // Remove existing notification
    const existingNotification = document.querySelector('.notification');
    if (existingNotification) {
        existingNotification.remove();
    }

    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.innerHTML = `
        <div class="notification-content">
            <span class="notification-message">${message}</span>
            <button class="notification-close" onclick="this.parentElement.parentElement.remove()">×</button>
        </div>
    `;
    
    // Add styles
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        min-width: 300px;
        max-width: 500px;
        background: ${getNotificationColor(type)};
        color: white;
        border-radius: 8px;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        z-index: 10000;
        animation: slideInRight 0.3s ease;
        font-family: inherit;
    `;

    document.body.appendChild(notification);

    // Auto remove after 5 seconds
    setTimeout(() => {
        if (notification.parentElement) {
            notification.remove();
        }
    }, 5000);
}

function getNotificationColor(type) {
    const colors = {
        success: '#2ed573',
        error: '#ff4757',
        warning: '#ffa502',
        info: '#3742fa'
    };
    return colors[type] || colors.info;
}

// Add notification styles to head
function injectNotificationStyles() {
    if (!document.querySelector('#notification-styles')) {
        const style = document.createElement('style');
        style.id = 'notification-styles';
        style.textContent = `
            .notification-content {
                padding: 12px 16px;
                display: flex;
                align-items: center;
                justify-content: between;
            }
            .notification-message {
                flex: 1;
                margin-right: 10px;
            }
            .notification-close {
                background: none;
                border: none;
                color: white;
                font-size: 18px;
                cursor: pointer;
                padding: 0;
                width: 24px;
                height: 24px;
                display: flex;
                align-items: center;
                justify-content: center;
                border-radius: 4px;
            }
            .notification-close:hover {
                background: rgba(255, 255, 255, 0.2);
            }
            @keyframes slideInRight {
                from {
                    transform: translateX(100%);
                    opacity: 0;
                }
                to {
                    transform: translateX(0);
                    opacity: 1;
                }
            }
        `;
        document.head.appendChild(style);
    }
}

// === PASSWORD STRENGTH CALCULATOR ===
function calculatePasswordStrength(password) {
    let score = 0;
    const checks = [
        password.length >= 8,
        /[a-z]/.test(password),
        /[A-Z]/.test(password),
        /[0-9]/.test(password),
        /[^a-zA-Z0-9]/.test(password)
    ];
    
    score = checks.filter(Boolean).length;

    const strengthLevels = [
        { percentage: 20, color: '#ff4757', text: 'Weak' },
        { percentage: 40, color: '#ffa502', text: 'Fair' },
        { percentage: 60, color: '#ffa502', text: 'Good' },
        { percentage: 80, color: '#2ed573', text: 'Strong' },
        { percentage: 100, color: '#2ed573', text: 'Very Strong' }
    ];
    
    return strengthLevels[Math.min(score, strengthLevels.length - 1)];
}

function validateForm(data, isSignup = false) {
    const errors = [];

    if (isSignup) {
        if (!data.name?.trim()) errors.push('Name is required');
        if (!data.user_type) errors.push('Please select user type');
    }

    if (!data.email?.trim()) {
        errors.push('Email is required');
    } else if (!isValidEmail(data.email)) {
        errors.push('Please enter a valid email address');
    }

    if (!data.password) {
        errors.push('Password is required');
    } else if (data.password.length < 6) {
        errors.push('Password must be at least 6 characters');
    }

    return errors;
}

// BETTER EMAIL VALIDATION
function isValidEmail(email) {
    // Simple check: has @, has ., and reasonable length
    if (!email.includes('@') || !email.includes('.') || email.length < 5) {
        return false;
    }
    return true
}

// === FORM HANDLERS ===
async function handleSignup(form) {
    const formData = new FormData(form);
    const data = {
        name: formData.get('name'),
        email: formData.get('email'),
        password: formData.get('password'),
        user_type: formData.get('userType')
    };

    // Validate form
    const errors = validateForm(data, true);
    if (errors.length > 0) {
        showNotification(errors[0], 'error');
        return;
    }

    // Show loading state
    const submitBtn = form.querySelector('button[type="submit"]');
    const originalText = submitBtn.textContent;
    submitBtn.innerHTML = '<div class="loading-spinner"></div> Creating Account...';
    submitBtn.disabled = true;

    try {
        const response = await fetch('/api/auth/signup', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        const result = await response.json();

        if (response.ok) {
            showNotification('🎉 Account created successfully! Redirecting...', 'success');
            
            // Store tokens and user data
            AuthManager.storeTokens(result.access_token, result.refresh_token, result.user);
            
            // Redirect to dashboard
            setTimeout(() => {
                window.location.href = '/dashboard';
            }, 2000);
        } else {
            showNotification(result.error || 'Signup failed. Please try again.', 'error');
        }
    } catch (error) {
        console.error('Signup error:', error);
        showNotification('🚫 Network error. Please check your connection.', 'error');
    } finally {
        submitBtn.textContent = originalText;
        submitBtn.disabled = false;
    }
}

async function handleLogin(form) {
    const formData = new FormData(form);
    const data = {
        email: formData.get('email'),
        password: formData.get('password')
    };

    // Validate form
    const errors = validateForm(data);
    if (errors.length > 0) {
        showNotification(errors[0], 'error');
        return;
    }

    // Show loading state
    const submitBtn = form.querySelector('button[type="submit"]');
    const originalText = submitBtn.textContent;
    submitBtn.innerHTML = '<div class="loading-spinner"></div> Signing In...';
    submitBtn.disabled = true;

    try {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        const result = await response.json();

        if (response.ok) {
            showNotification('👋 Welcome back! Redirecting...', 'success');
            
            // Store tokens and user data
            AuthManager.storeTokens(result.access_token, result.refresh_token, result.user);
            
            // Redirect to dashboard
            setTimeout(() => {
                window.location.href = '/dashboard';
            }, 1500);
        } else {
            showNotification(result.error || 'Login failed. Please check your credentials.', 'error');
        }
    } catch (error) {
        console.error('Login error:', error);
        showNotification('🚫 Network error. Please try again.', 'error');
    } finally {
        submitBtn.textContent = originalText;
        submitBtn.disabled = false;
    }
}

function handleSocialAuth(provider) {
    showNotification(`🔜 ${provider} authentication coming soon!`, 'info');
}

// === UI CONTROLS ===
function setupFormToggles() {
    const toggleBtns = document.querySelectorAll('.toggle-btn');
    const authForms = document.querySelectorAll('.auth-form');
    const switchLinks = document.querySelectorAll('.switch-link');
    
    // Toggle between signup and login
    toggleBtns.forEach(btn => {
        btn.addEventListener('click', function() {
            const formType = this.getAttribute('data-form');
            
            // Update active toggle button
            toggleBtns.forEach(b => b.classList.remove('active'));
            this.classList.add('active');
            
            // Show corresponding form
            authForms.forEach(form => {
                form.classList.remove('active');
                if (form.id === `${formType}-form`) {
                    form.classList.add('active');
                }
            });
        });
    });
    
    // Switch links
    switchLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const formType = this.getAttribute('data-form');
            const targetBtn = document.querySelector(`.toggle-btn[data-form="${formType}"]`);
            if (targetBtn) targetBtn.click();
        });
    });
}

function setupUserTypeToggle() {
    const userTypeRadios = document.querySelectorAll('input[name="userType"]');
    const trainerFields = document.querySelector('.trainer-fields');
    
    userTypeRadios.forEach(radio => {
        radio.addEventListener('change', function() {
            trainerFields.style.display = this.value === 'trainer' ? 'block' : 'none';
        });
    });
}

function setupPasswordStrength() {
    const passwordInput = document.getElementById('signup-password');
    const strengthBar = document.querySelector('.strength-bar');
    const strengthText = document.querySelector('.strength-text');
    
    if (passwordInput && strengthBar && strengthText) {
        passwordInput.addEventListener('input', function() {
            const password = this.value;
            const strength = calculatePasswordStrength(password);
            
            strengthBar.style.width = `${strength.percentage}%`;
            strengthBar.style.background = strength.color;
            strengthText.textContent = strength.text;
            strengthText.style.color = strength.color;
        });
    }
}

function setupFormSubmissions() {
    const signupForm = document.getElementById('signup-form');
    const loginForm = document.getElementById('login-form');
    
    if (signupForm) {
        signupForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            await handleSignup(signupForm);
        });
    }
    
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            await handleLogin(loginForm);
        });
    }
}

function setupSocialAuth() {
    const socialButtons = document.querySelectorAll('.btn-google, .btn-apple');
    socialButtons.forEach(btn => {
        btn.addEventListener('click', function() {
            const provider = this.classList.contains('btn-google') ? 'Google' : 'Apple';
            handleSocialAuth(provider);
        });
    });
}

// === INITIALIZATION ===
document.addEventListener('DOMContentLoaded', function() {
    // Inject styles
    injectNotificationStyles();
    
    // Check if user is already logged in
    if (AuthManager.isAuthenticated() && window.location.pathname === '/auth') {
        window.location.href = '/dashboard';
        return;
    }
    
    // Setup all UI controls
    setupFormToggles();
    setupUserTypeToggle();
    setupPasswordStrength();
    setupFormSubmissions();
    setupSocialAuth();
    
    console.log('🚀 FitterBy Auth Frontend Loaded');
});

// Add loading spinner styles
const spinnerStyles = document.createElement('style');
spinnerStyles.textContent = `
    .loading-spinner {
        display: inline-block;
        width: 16px;
        height: 16px;
        border: 2px solid #ffffff;
        border-radius: 50%;
        border-top-color: transparent;
        animation: spin 1s ease-in-out infinite;
        margin-right: 8px;
    }
    @keyframes spin {
        to { transform: rotate(360deg); }
    }
`;
document.head.appendChild(spinnerStyles);
