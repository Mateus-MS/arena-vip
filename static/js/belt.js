import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';

const canvas = document.getElementById('belt-canvas');
if (!canvas) throw new Error('belt-canvas not found');

canvas.style.opacity = '0';

const isMobile = window.innerWidth < 1024;
canvas.parentElement.style.height = isMobile ? '20dvh' : 'clamp(500px, 80vh, 900px)';

// ── Start the GLB download immediately so it's in-flight (or done) by the
// time the observer fires. Three.js setup is still deferred to first visibility.
const modelUrl = isMobile ? '/static/belt_mobile.glb' : '/static/belt.glb';
const modelReady = new Promise((resolve, reject) => {
    new GLTFLoader().load(modelUrl, resolve, undefined, reject);
});

// ── Three.js init — called by the observer, consumes the pre-downloaded model ─
function initBelt() {
    const renderer = new THREE.WebGLRenderer({ canvas, alpha: true, antialias: true });
    const dpr = Math.min(window.devicePixelRatio, 2);
    renderer.setPixelRatio(dpr);
    renderer.outputColorSpace = THREE.SRGBColorSpace;
    renderer.toneMapping = THREE.ACESFilmicToneMapping;
    renderer.toneMappingExposure = 1.0;

    const scene  = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(28, 1, 0.01, 200);
    if (isMobile) { camera.position.set(0, 0, 4.5); }
    else          { camera.position.set(0, 0.3, 7.5); }
    camera.lookAt(0, 0, 0);

    scene.add(new THREE.AmbientLight(0xffffff, 0.05));
    const key  = new THREE.DirectionalLight(0xffffff, 1.8); key.position.set(2, 5, 4);    scene.add(key);
    const fill = new THREE.DirectionalLight(0xffffff, 0.25); fill.position.set(-4, 1, 2); scene.add(fill);
    const rim  = new THREE.DirectionalLight(0xffffff, 3.0);  rim.position.set(-1, -2, -5); scene.add(rim);
    const rim2 = new THREE.DirectionalLight(0xffffff, 1.2);  rim2.position.set(1, 4, -4);  scene.add(rim2);

    const uTime       = { value: 0 };
    const uResolution = { value: new THREE.Vector2(1, 1) };
    const uDissolve   = { value: 1.0 };

    const NOISE_GLSL = `
float _h(vec2 p){p=fract(p*vec2(127.1,311.7));p+=dot(p,p+45.32);return fract(p.x*p.y);}
float _vn(vec2 p){
    vec2 i=floor(p),f=fract(p);f=f*f*(3.0-2.0*f);
    return mix(mix(_h(i),_h(i+vec2(1,0)),f.x),mix(_h(i+vec2(0,1)),_h(i+vec2(1,1)),f.x),f.y);}
float _fbm(vec2 p){float v=0.,a=.5;for(int i=0;i<5;i++){v+=a*_vn(p);p*=2.1;a*=.5;}return v;}
`;

    function injectShaders(mat, localYMax) {
        const uYMax = { value: localYMax };
        mat.onBeforeCompile = (shader) => {
            shader.uniforms.uTime       = uTime;
            shader.uniforms.uYMax       = uYMax;
            shader.uniforms.uResolution = uResolution;
            shader.uniforms.uDissolve   = uDissolve;

            shader.vertexShader = `uniform float uTime;\nuniform float uYMax;\n` + shader.vertexShader;
            shader.vertexShader = shader.vertexShader.replace('#include <begin_vertex>',
                `#include <begin_vertex>
                float tipFactor = pow(clamp(position.y/uYMax,0.,1.),2.);
                float wave = sin(uTime*1.2+position.x*4.+position.y*1.5)*.55
                           + sin(uTime*0.7+position.x*6.5-position.y*1.)*.45;
                transformed.z += wave*.03*tipFactor;
                transformed.x += wave*.01*tipFactor;`
            );

            shader.fragmentShader =
                `uniform vec2 uResolution;\nuniform float uTime;\nuniform float uDissolve;\n`
                + NOISE_GLSL + shader.fragmentShader;

            shader.fragmentShader = shader.fragmentShader.replace('#include <normal_fragment_begin>',
                `{
                    vec2 noiseCoord = gl_FragCoord.xy/uResolution.x*10.;
                    if(_fbm(noiseCoord)<uDissolve) discard;
                    vec2 sUV=gl_FragCoord.xy/uResolution;
                    float wv=sin(sUV.y*3.5+uTime*.8)*.10+sin(sUV.y*8.-uTime*1.3)*.04;
                    float t=step(0.,(sUV.x-sUV.y)+wv);
                    bool isRed=(diffuseColor.r-max(diffuseColor.g,diffuseColor.b))>.005;
                    diffuseColor.rgb=mix(diffuseColor.rgb,isRed?vec3(0):vec3(1),t);
                }
                #include <normal_fragment_begin>`
            );

            shader.fragmentShader = shader.fragmentShader.replace('#include <emissivemap_fragment>',
                `#include <emissivemap_fragment>
                {
                    vec2 sUV=gl_FragCoord.xy/uResolution;
                    float wv=sin(sUV.y*3.5+uTime*.8)*.10+sin(sUV.y*8.-uTime*1.3)*.04;
                    float absD=abs((sUV.x-sUV.y)+wv);
                    float core =1.-smoothstep(.0000,.0015,absD);
                    float inner=1.-smoothstep(.0015,.0060,absD);
                    float outer=1.-smoothstep(.0060,.0180,absD);
                    float pulse=.88+.12*sin(uTime*1.8);
                    totalEmissiveRadiance+=core*vec3(12,11,9)*pulse+inner*vec3(3.5,2.5,.7)*pulse+outer*vec3(.8,.45,.04);
                    vec2 nc=gl_FragCoord.xy/uResolution.x*10.;
                    float n=_fbm(nc),edge=n-uDissolve;
                    float crackle=.85+.15*sin(uTime*18.+n*40.);
                    float gf=smoothstep(0.,.10,uDissolve);
                    totalEmissiveRadiance+=(1.-smoothstep(.000,.010,edge))*gf*crackle*vec3(24,20,8)
                        +(1.-smoothstep(.010,.048,edge))*gf*vec3(7,4.5,.6)
                        +(1.-smoothstep(.048,.110,edge))*gf*vec3(2,.9,.05);
                }`
            );
        };
        mat.customProgramCacheKey = () => `belt-${localYMax.toFixed(5)}`;
        mat.needsUpdate = true;
    }

    let belt     = null;
    let loadedAt = null;
    let prevMs   = 0;

    // Consume the already-in-flight (or resolved) download promise
    modelReady.then((gltf) => {
        const obj    = gltf.scene;
        const box    = new THREE.Box3().setFromObject(obj);
        const center = box.getCenter(new THREE.Vector3());
        const size   = box.getSize(new THREE.Vector3());

        obj.position.sub(center);

        const pivot = new THREE.Group();
        pivot.add(obj);

        if (isMobile) {
            pivot.scale.setScalar(5.47 / Math.max(size.x, size.y, size.z));
            pivot.position.set(0.6, 0, 0);
            pivot.rotation.x = -0.08;
        } else {
            pivot.scale.setScalar(2.64 / Math.max(size.x, size.y, size.z));
            pivot.position.y = -0.1;
            pivot.rotation.x = -0.18;
        }

        scene.add(pivot);
        belt = pivot;

        obj.traverse((child) => {
            if (!child.isMesh) return;
            const lb        = new THREE.Box3().setFromBufferAttribute(child.geometry.attributes.position);
            const localYMax = lb.max.y > 0.001 ? lb.max.y * 0.80 : 9999.0;
            const mats      = Array.isArray(child.material) ? child.material : [child.material];
            mats.forEach((m, i) => {
                const cloned = m.clone();
                if (Array.isArray(child.material)) child.material[i] = cloned;
                else child.material = cloned;
                injectShaders(cloned, localYMax);
            });
        });

        loadedAt = performance.now();
        canvas.style.opacity = '1';
    });

    function syncSize() {
        const w = canvas.clientWidth;
        const h = canvas.clientHeight;
        if (renderer.domElement.width !== w * dpr || renderer.domElement.height !== h * dpr) {
            renderer.setSize(w, h, false);
            camera.aspect = w / h;
            camera.updateProjectionMatrix();
        }
        uResolution.value.set(w * dpr, h * dpr);
    }

    function animate(ms) {
        requestAnimationFrame(animate);
        syncSize();

        prevMs = ms;
        uTime.value = ms * 0.001;

        if (belt) {
            const t = ms * 0.001;

            if (loadedAt !== null) {
                const p = Math.min((ms - loadedAt) / 2400, 1);
                uDissolve.value = 1.0 - (1 - Math.pow(1 - p, 3));
            }

            belt.rotation.y = Math.sin(t * 0.25) * (isMobile ? 0.06 : 0.12);
            belt.position.y = (isMobile ? 0 : -0.1) + Math.sin(t * 0.55) * (isMobile ? 0.03 : 0.05);
        }

        renderer.render(scene, camera);
    }

    animate(0);
}

// ── Observer: triggers Three.js init on first intersection ───────────────────
const observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting) {
        observer.disconnect();
        initBelt();
    }
}, { threshold: 0 });

observer.observe(canvas);
