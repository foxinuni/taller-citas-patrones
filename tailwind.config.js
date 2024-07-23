/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./core/views/**/*.templ"],
  theme: {
    extend: {
      colors: {
        "text": "#fafafa",
        "back": "#18181b",
        "primary": "#e11d48",
        "primary-light": "#ff9b54",
        "secondary": "#3f3f46",
        "accent": "#be123c",
      }
    },
  },
  plugins: [],
}

