/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx}'
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif']
      },
      container: {
        center: true,
        padding: '1rem'
      },
      colors: {
        brand: {
          primary: '#2563EB',
          accent: '#059669'
        }
      }
    }
  },
  plugins: []
}
