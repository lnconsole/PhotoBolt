/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        pbrown: {
          500: '#20acc7'
        },
        pbrown2: {
          500: '#1d9bb3'
        },
        pbrown3: {
          500: '#1d9bb3'
        },
        pbrown4: {
          500: '#20acc7'
        },
        pbg: {
          500: '#f2f2f2'
        },
        pbg2: {
          500: '#a6dee9'
        },
        ptxtl: {
          500: '#f8dfdc'
        },
        ptxtd: {
          500: '#555555'
        },
        btnl: {
          500: '#f8dfdc'
        },
        btntxtl: {
          500: '#2e2e2e'
        },
        btna: {
          500: '#8c0327'
        },
        btntxta: {
          500: '#edd0cf'
        },
      }
    },
  },
  plugins: [],
}

