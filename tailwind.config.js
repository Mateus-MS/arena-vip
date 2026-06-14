/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './**/*.templ',
    './**/*.go',
    '!./node_modules/**',
    '!./bin/**',
  ],
  theme: {
    extend: {
      colors: {
        /* ── Semantic theme tokens ─────────────────────────────────
           All mapped to CSS variables in static/css/input.css.
           Change the variables there → whole site recolors.        */
        canvas:   'rgb(var(--canvas)   / <alpha-value>)',
        surface:  'rgb(var(--surface)  / <alpha-value>)',
        contrast: 'rgb(var(--contrast) / <alpha-value>)',
        ink:      'rgb(var(--ink)      / <alpha-value>)',
        muted:    'rgb(var(--muted)    / <alpha-value>)',
        subtle:   'rgb(var(--subtle)   / <alpha-value>)',
        edge:     'rgb(var(--edge)     / <alpha-value>)',
        accent:   'rgb(var(--accent)   / <alpha-value>)',
        gold:     'rgb(var(--gold)     / <alpha-value>)',
      },
      fontFamily: {
        display: ['"Bebas Neue"', 'sans-serif'],
        body:    ['Inter', 'sans-serif'],
        accent:  ['Oswald', 'sans-serif'],
      },
      animation: {
        'spin-slow':    'spin 22s linear infinite',
        'spin-reverse': 'spin-reverse 16s linear infinite',
        'float':        'float 7s ease-in-out infinite',
        'fade-in-up':   'fadeInUp 0.7s ease-out forwards',
      },
      keyframes: {
        'spin-reverse': {
          from: { transform: 'rotate(360deg)' },
          to:   { transform: 'rotate(0deg)' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%':      { transform: 'translateY(-12px)' },
        },
        fadeInUp: {
          from: { opacity: '0', transform: 'translateY(28px)' },
          to:   { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}
