# Arena VIP

> Premium martial arts academy website — Santo Tirso, Portugal.

A production-ready web application for a high-end combat sports academy offering Brazilian Jiu-Jitsu, Kickboxing, and Judo. The site is built with a philosophy of restraint: no JavaScript framework on the server, no unnecessary abstractions, and a visual identity that earns the word *premium* — dark, typographic, and uncompromising.

---

## Concept

Arena VIP is not a sports club. It is positioned as a *discipline forge* — a place where the practitioner is built from the inside out. The brand language reflects this: near-black backgrounds, a restricted palette of white, crimson, and gold, and display typography at aggressive scale. Every page is designed to feel like an institution, not a landing page.

The centerpiece of the homepage is a real-time 3D render of a black BJJ belt — the highest symbol of commitment in the discipline — dissolving into existence through a GPU noise shader and split by an animated gold seam that glows in the dark.

---

## Stack

| Layer | Technology |
|---|---|
| Server | Go 1.23 + [Gin](https://github.com/gin-gonic/gin) v1.9 |
| Templating | [Templ](https://templ.guide/) v0.3 (compiled to Go) |
| Styling | Tailwind CSS v3 (JIT) |
| Interactivity | Alpine.js v3 (via CDN, `defer`) |
| 3D Rendering | Three.js r170 (ES module via importmap) |
| Email | Go `net/smtp` → Gmail SMTP (port 587, STARTTLS) |
| Fonts | Bebas Neue · Oswald · Inter (Google Fonts) |

No Node.js runtime is required at runtime. Tailwind runs at build time; the compiled CSS is committed to `static/css/output.css`.

---

## Project Structure

```
arena-vip/
├── main.go                     # Entry point: loads .env, mounts static, starts :8080
├── go.mod
│
├── app/
│   ├── app.go                  # Singleton Gin engine + security middleware
│   └── config/config.go        # Feature flags + env var accessors
│
├── templates/
│   └── base.templ              # HTML shell: <head>, OG tags, JSON-LD, importmap
│
├── atoms/                      # Primitive UI components (Button, Input, Badge, …)
├── components/                 # Composite UI components (Navbar, Footer, Contact form, …)
│
├── pages/
│   ├── landing/                # / — Homepage
│   │   └── frags/              # Hero, Modalities, About, Location, Join, Social proof
│   ├── horarios/               # /horarios — Weekly class schedule
│   ├── professores/            # /professores — Coach profiles
│   ├── resultados/             # /resultados — Competition results
│   ├── depoimentos/            # /depoimentos — Student testimonials
│   ├── faq/                    # /faq — Frequently asked questions
│   ├── loja/                   # /loja — Merch shop (feature-flagged)
│   ├── contato/                # POST /contato — Contact form handler + rate limiter
│   └── privacidade/            # /privacidade — Privacy policy
│
├── utils/
│   └── render.go               # Templ → Gin response helper
│
└── static/
    ├── belt.glb                # Full-resolution 3D belt model (~32 MB)
    ├── belt_mobile.glb         # Mobile-optimised variant
    ├── favicon.svg
    ├── js/
    │   └── belt.js             # Three.js + GLSL shader controller
    └── css/
        ├── input.css           # Theme tokens + Tailwind directives
        └── output.css          # Compiled output (committed)
```

### Page Registration Pattern

Each page package self-registers its route inside `init()`. `main.go` blank-imports every package to trigger registration — there is no central route file.

```go
// pages/horarios/page.templ
func init() {
    app.GetInstance().Router.GET("/horarios", utils.Render(Page()))
}
```

The `app.GetInstance()` singleton is protected by `sync.Once`, so concurrent `init()` calls across packages are safe.

---

## Design System

All colour decisions live in a single block of CSS custom properties using **bare RGB channels** so Tailwind's opacity modifiers (`text-accent/50`, `bg-gold/10`) work without extra configuration.

```css
:root {
  --canvas:   255 255 255;   /* page background — white          */
  --contrast: 8   8   8;     /* dark section bg — near-black     */
  --ink:      10  10  10;    /* primary text                     */
  --muted:    107 114 128;   /* secondary text                   */
  --edge:     226 232 240;   /* border colour                    */
  --accent:   190 18  40;    /* crimson — BJJ belt tip           */
  --gold:     201 168 76;    /* gold — premium accent            */
}
```

**Typography roles:**

| Class | Font | Usage |
|---|---|---|
| `font-display` | Bebas Neue | Hero headlines, stats, grid headers |
| `font-accent` | Oswald | Labels, badges, nav links, uppercase UI text |
| `font-body` | Inter | Body copy, descriptions, form labels |

---

## The 3D Belt

The homepage hero renders a real black belt GLB model in a `<canvas>` using Three.js. The entire system is split into two phases to balance load performance against visual impact.

### Phase 1 — Prefetch on page load

The GLB download starts the moment the JS module is parsed, before any Three.js renderer exists:

```js
const modelReady = new Promise((resolve, reject) => {
    new GLTFLoader().load(modelUrl, resolve, undefined, reject);
});
```

Mobile gets `belt_mobile.glb`; desktop gets the full `belt.glb`.

### Phase 2 — Three.js init on intersection

The renderer, scene, camera, and lights are only created when the canvas enters the viewport, via `IntersectionObserver`. This avoids spending GPU resources on a canvas the user has not yet seen:

```js
const observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting) {
        observer.disconnect();
        initBelt(); // consumes the already-in-flight modelReady promise
    }
}, { threshold: 0 });
observer.observe(canvas);
```

### GLSL Shader Injection

Standard Three.js `MeshStandardMaterial` is extended at shader compile time via `onBeforeCompile`. Two visual effects are injected:

**Dissolve entrance** — The belt materialises out of noise over 2.4 seconds. An FBM (Fractional Brownian Motion) function generates a smooth noise field in screen space. Fragments below a `uDissolve` threshold are `discard`ed; fragments near the dissolve boundary emit a hot-gold crackle glow via `totalEmissiveRadiance`.

```glsl
float _fbm(vec2 p) {
    float v = 0., a = .5;
    for (int i = 0; i < 5; i++) { v += a * _vn(p); p *= 2.1; a *= .5; }
    return v;
}
// ...
if (_fbm(noiseCoord) < uDissolve) discard;
```

The dissolve factor is driven by a cubic ease-out on the CPU:
```js
const p = Math.min((ms - loadedAt) / 2400, 1);
uDissolve.value = 1.0 - (1 - Math.pow(1 - p, 3));
```

**Belt seam** — A diagonal split in screen UV space (`sUV.x - sUV.y`) divides the belt into a black half and a white half, recreating the visual of a folded belt. A compound sine wave distorts the boundary over time, making the seam ripple. A multi-layered emissive glow (core → inner → outer) traces the seam boundary in gold.

---

## Class Schedule Page (`/horarios`)

The schedule page is fully server-rendered with dynamic "today" detection using Lisbon local time:

```go
var lisboaLoc = func() *time.Location {
    loc, err := time.LoadLocation("Europe/Lisbon")
    if err != nil { return time.UTC }
    return loc
}()
```

The page is divided into three sections:

1. **Hero** — Page title + two quick-fact stats (training days per week, number of disciplines).
2. **Today Spotlight** — Shows the current day's classes as large cards with time, modality, and instructor. On rest days, shows the next training day and its class pills.
3. **Week Grid** — A 6-column grid (Mon–Sat) that highlights today in gold and dims inactive days. On mobile it collapses to 2 columns.

---

## Contact Form & Security

### Spam prevention

- **Honeypot field** — A hidden `<input name="website">` is positioned off-screen. Bots fill it; humans don't. On a non-empty value the server returns `200 OK` with a fake success — no signal to the bot.
- **Rate limiting** — In-memory sliding-window limiter (3 requests / 10 minutes per IP), protected by `sync.Mutex`.

### Input validation

- Request body capped at **64 KB** via `http.MaxBytesReader` before form parsing begins.
- Required fields, length limits (name ≤ 100, email ≤ 200, message ≤ 5000), email format check, and a `modalidade` allowlist are all validated server-side.

### HTTP security headers

Applied globally via a Gin middleware:

```
X-Frame-Options: SAMEORIGIN
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
```

### Cache headers

```
/static/*  →  Cache-Control: public, max-age=31536000, immutable
pages      →  Cache-Control: public, max-age=300, must-revalidate
```

---

## SEO

**`base.templ`** emits a full metadata suite on every page:

- `<link rel="canonical">` — canonical URL
- Open Graph tags — title, description, type, URL, image (1200×630), locale, site name
- Twitter Card — `summary_large_image`
- **JSON-LD** — `SportsActivityLocation` schema with address, coordinates, phone, email, opening hours, and social profile links

---

## Environment Variables

Create a `.env` file in the project root. The server reads it at startup via a minimal parser in `main.go` (no third-party dotenv dependency).

| Variable | Required | Description |
|---|---|---|
| `SMTP_USER` | Yes (for email) | Gmail address used as sender |
| `SMTP_PASS` | Yes (for email) | Gmail app password |
| `CONTACT_TO` | Yes (for email) | Recipient address for contact form submissions |
| `MAPS_API_KEY` | No | Google Maps JavaScript API key |
| `MAPS_ADDRESS` | No | Gym address for the map pin |
| `SHOP_ENABLED` | No | Set to `true` to show the `/loja` route and nav link |

---

## Development

**Prerequisites:** Go 1.23+, Templ CLI, Node.js (for Tailwind only).

```bash
# Install Templ CLI
go install github.com/a-h/templ/cmd/templ@latest

# Generate Go code from .templ files
templ generate ./...

# Compile Tailwind CSS
npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

# Run the server
go run .
# → http://localhost:8080
```

**Build order matters:** `.templ` files must be compiled to `_templ.go` before `go build` or `go run` will succeed. Always run `templ generate` after editing any `.templ` file.

```bash
# Full build
templ generate ./... && go build ./...
```

---

## Deployment

The application is a single self-contained Go binary with no runtime dependencies beyond the operating system. Static assets are served by Gin's `r.Static()` handler.

```bash
templ generate ./...
go build -o arena-vip .
./arena-vip
```

Serve behind a reverse proxy (Caddy, nginx, or a cloud load balancer) that terminates TLS. The binary listens on `:8080`.

---

*© Arena VIP — Santo Tirso, Portugal*
