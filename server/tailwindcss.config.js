const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/ui/views/**/*.html"],
  plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
};
