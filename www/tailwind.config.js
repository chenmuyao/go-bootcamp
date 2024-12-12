/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        'gumbo': {
          '50': '#f3f8f8',
          '100': '#e1eaec',
          '200': '#c6d8db',
          '300': '#9ebcc2',
          '400': '#82a5ad',
          '500': '#547c86',
          '600': '#486872',
          '700': '#3f565f',
          '800': '#394a51',
          '900': '#334046',
          '950': '#1f282d',
        },
      },
    },
  },
  plugins: [],
}

