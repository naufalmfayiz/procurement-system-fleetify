/**
 * Reusable AJAX API Module
 * Handles all API requests with automatic token injection
 */

// Token management functions
function getToken() {
  return localStorage.getItem("jwt_token");
}

function setToken(token) {
  localStorage.setItem("jwt_token", token);
}

function removeToken() {
  localStorage.removeItem("jwt_token");
}

function getUser() {
  const user = localStorage.getItem("user");
  return user ? JSON.parse(user) : null;
}

function setUser(user) {
  localStorage.setItem("user", JSON.stringify(user));
}

function removeUser() {
  localStorage.removeItem("user");
}

function logout() {
  removeToken();
  removeUser();
  window.location.href = "index.html";
}

// Check if user is authenticated
function requireAuth() {
  if (!getToken()) {
    window.location.href = "index.html";
    return false;
  }
  return true;
}

/**
 * Reusable AJAX wrapper
 * Automatically adds Authorization header with JWT token
 */
const api = {
  /**
   * Make a GET request
   * @param {string} endpoint - API endpoint (without base URL)
   * @returns {jqXHR} jQuery AJAX promise
   */
  get: function (endpoint) {
    return this.request("GET", endpoint);
  },

  /**
   * Make a POST request
   * @param {string} endpoint - API endpoint (without base URL)
   * @param {object} data - Request body data
   * @returns {jqXHR} jQuery AJAX promise
   */
  post: function (endpoint, data) {
    return this.request("POST", endpoint, data);
  },

  /**
   * Make a PUT request
   * @param {string} endpoint - API endpoint (without base URL)
   * @param {object} data - Request body data
   * @returns {jqXHR} jQuery AJAX promise
   */
  put: function (endpoint, data) {
    return this.request("PUT", endpoint, data);
  },

  /**
   * Make a DELETE request
   * @param {string} endpoint - API endpoint (without base URL)
   * @returns {jqXHR} jQuery AJAX promise
   */
  delete: function (endpoint) {
    return this.request("DELETE", endpoint);
  },

  /**
   * Core AJAX request function
   * @param {string} method - HTTP method
   * @param {string} endpoint - API endpoint
   * @param {object} data - Optional request body
   * @returns {jqXHR} jQuery AJAX promise
   */
  request: function (method, endpoint, data = null) {
    const url = API_BASE_URL + endpoint;
    const token = getToken();

    const options = {
      url: url,
      method: method,
      contentType: "application/json",
      dataType: "json",
    };

    // Add Authorization header if token exists
    if (token) {
      options.headers = {
        Authorization: "Bearer " + token,
      };
    }

    // Add request body for POST/PUT
    if (data && (method === "POST" || method === "PUT")) {
      options.data = JSON.stringify(data);
    }

    // Create AJAX request with error handling
    return $.ajax(options).fail(function (xhr) {
      // Handle 401 Unauthorized - redirect to login
      if (xhr.status === 401) {
        toastr.error("Session expired. Please login again.");
        logout();
      }
    });
  },
};

// Format currency (IDR)
function formatCurrency(amount) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(amount);
}

// Format date
function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleDateString("id-ID", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}
