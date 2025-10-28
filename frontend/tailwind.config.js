/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        dark: {
          bg: '#1e272e',
          card: '#2d3436',
          lighter: '#3d4547',
        },
        accent: {
          cyan: '#00d2d3',
          purple: '#a29bfe',
          pink: '#fd79a8',
        },
      },
      fontFamily: {
        sans: ['Poppins', 'system-ui', '-apple-system', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
