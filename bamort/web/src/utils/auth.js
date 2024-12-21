// src/utils/auth.js
export function isLoggedIn() {
  const token = localStorage.getItem("token");
  return !!token; // Returns true if the token exists
}

export function logout() {
  localStorage.removeItem("token");
}
