<script setup lang="ts">
const { t } = useI18n()
const localePath = useLocalePath()
</script>

<template>
  <section :class="$style.hero">
    <!-- Animated Background Elements -->
    <div :class="$style.heroBgWrapper">
      <div :class="$style.heroGlow" />
      <div :class="$style.heroGradientBg" />
      <div :class="$style.heroGridBg" />
      <div :class="$style.heroAnimatedShapes">
        <div :class="[$style.shape, $style.shape1]" />
        <div :class="[$style.shape, $style.shape2]" />
        <div :class="[$style.shape, $style.shape3]" />
      </div>
      <div :class="$style.heroCodeSnippets">
        <div :class="[$style.codeSnippet, $style.snippet1]">
          <div :class="$style.codeLine">
            <span :class="$style.codeKeyword">func</span>
            <span :class="$style.codeFunction">GetPost</span>() {'{'}
          </div>
          <div :class="$style.codeLine">
            <span :class="$style.codeKeyword">return</span> post
          </div>
          <div :class="$style.codeLine">{'}'}</div>
        </div>
        <div :class="[$style.codeSnippet, $style.snippet2]">
          <div :class="$style.codeLine">
            <span :class="$style.codeTag">&lt;template&gt;</span>
          </div>
          <div :class="$style.codeLine">
            <span :class="$style.codeTag">&lt;div</span>
            <span :class="$style.codeAttr">class</span>=
            <span :class="$style.codeString">"content-hub"</span>
            <span :class="$style.codeTag">&gt;</span>
          </div>
          <div :class="$style.codeLine">
            <span :class="$style.codeTag">&lt;/div&gt;</span>
          </div>
        </div>
        <div :class="[$style.codeSnippet, $style.snippet3]">
          <div :class="$style.codeLine">
            <span :class="$style.codeComment">// Content Hub API</span>
          </div>
          <div :class="$style.codeLine">
            <span :class="$style.codeKeyword">const</span> api =
            <span :class="$style.codeString">'v1'</span>
          </div>
        </div>
      </div>
      <svg :class="$style.heroWaves" viewBox="0 0 1440 320" preserveAspectRatio="none">
        <path fill="hsl(var(--primary) / 0.08)" d="M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,122.7C672,117,768,139,864,138.7C960,139,1056,117,1152,101.3C1248,85,1344,75,1392,69.3L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z" />
      </svg>
    </div>

    <div :class="$style.heroContent">
      <h1 :class="$style.heroTitle">{{ t('authentication.login.brand_title') }}</h1>
      <p :class="$style.heroSubtitle">{{ t('authentication.login.brand_subtitle') }}</p>
      <p :class="$style.heroDescription">{{ t('page.home.hero_description') }}</p>
      <div :class="$style.heroActions">
        <UiButton size="lg" :class="$style.btnPrimary" @click="navigateTo(localePath('/post'))">
          {{ t('page.home.browse_posts') }}
        </UiButton>
        <UiButton variant="outline" size="lg" :class="$style.btnSecondary" @click="navigateTo(localePath('/about'))">
          {{ t('page.home.learn_more') }}
        </UiButton>
      </div>
    </div>
  </section>
</template>

<style module>
.hero {
    position: relative;
    padding: 7rem 2rem 2.5rem;
    text-align: center;
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 45%, #1e1b4b 100%);
    color: hsl(0 0% 98%);
    overflow: hidden;
    min-height: 320px;
    display: flex;
    align-items: center;
    justify-content: center;
}

:global(.dark) .hero {
    background: linear-gradient(45deg, #000000 0%, #0f172a 50%, rgba(6, 78, 59, 0.35) 100%);
}

.hero::after {
    content: '';
    position: absolute;
    bottom: 0; left: 0; right: 0;
    height: 100px;
    background: linear-gradient(to bottom, transparent 0%, hsl(var(--background)) 100%);
    pointer-events: none;
    z-index: 1;
}

.heroBgWrapper {
    position: absolute;
    inset: 0;
    overflow: hidden;
    z-index: 0;
    pointer-events: none;
}

.heroGlow {
    position: absolute;
    top: 50%; left: 50%;
    transform: translate(-50%, -50%);
    width: 400px; height: 400px;
    border-radius: 50%;
    background: hsl(var(--primary) / 0.08);
    filter: blur(100px);
    animation: glowBreathe 6s ease-in-out infinite;
}

@keyframes glowBreathe {
    0%, 100% { opacity: 0.6; transform: translate(-50%, -50%) scale(1); }
    50% { opacity: 1; transform: translate(-50%, -50%) scale(1.15); }
}

.heroGradientBg {
    position: absolute; inset: 0;
    background:
        radial-gradient(ellipse at 20% 30%, hsl(var(--primary) / 0.18) 0%, transparent 55%),
        radial-gradient(ellipse at 80% 70%, rgba(56, 189, 248, 0.12) 0%, transparent 55%),
        radial-gradient(ellipse at 50% 100%, hsl(var(--primary) / 0.15) 0%, transparent 60%);
    animation: gradientShift 15s ease-in-out infinite;
}

@keyframes gradientShift {
    0%, 100% { opacity: 1; transform: scale(1); }
    50% { opacity: 0.85; transform: scale(1.05); }
}

.heroGridBg {
    position: absolute; inset: 0;
    background-image:
        linear-gradient(hsl(var(--primary) / 0.06) 1px, transparent 1px),
        linear-gradient(90deg, hsl(var(--primary) / 0.06) 1px, transparent 1px);
    background-size: 50px 50px;
    animation: gridMove 20s linear infinite;
    opacity: 0.8;
    mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
    -webkit-mask-image: radial-gradient(ellipse at center, black 30%, transparent 80%);
}

@keyframes gridMove {
    0% { transform: translate(0, 0); }
    100% { transform: translate(50px, 50px); }
}

.heroAnimatedShapes { position: absolute; width: 100%; height: 100%; }

.shape { position: absolute; border-radius: 50%; opacity: 0.2; filter: blur(80px); }
.shape1 { width: 350px; height: 350px; background: hsl(var(--primary)); top: -150px; right: -100px; animation: float 8s ease-in-out infinite; }
.shape2 { width: 450px; height: 450px; background: hsl(var(--primary)); bottom: -200px; left: -150px; animation: float 10s ease-in-out infinite reverse; }
.shape3 { width: 250px; height: 250px; background: hsl(var(--primary)); top: 50%; left: 50%; animation: float 12s ease-in-out infinite; }

@keyframes float {
    0%, 100% { transform: translateY(0) scale(1); }
    50% { transform: translateY(-30px) scale(1.1); }
}

.heroCodeSnippets { position: absolute; width: 100%; height: 100%; pointer-events: none; perspective: 1000px; opacity: 0.12; }

.codeSnippet {
    position: absolute;
    background: hsl(var(--card) / 0.8);
    border: 1px solid hsl(var(--primary) / 0.2);
    border-radius: 8px;
    padding: 24px 36px;
    font-family: 'Fira Code', 'JetBrains Mono', 'Consolas', monospace;
    font-size: 28px; line-height: 1.6;
    box-shadow: 0 0 20px hsl(var(--primary) / 0.15);
    backdrop-filter: blur(4px);
}

.codeLine { white-space: nowrap; color: hsl(0 0% 85%); text-shadow: 0 0 8px hsl(var(--primary) / 0.3); }
.codeKeyword { color: hsl(var(--primary)); font-weight: 600; }
.codeFunction { color: hsl(48 96% 60%); }
.codeTag { color: hsl(200 80% 65%); }
.codeAttr { color: hsl(280 70% 70%); }
.codeString { color: hsl(var(--primary)); }
.codeComment { color: hsl(0 0% 45%); font-style: italic; }

.snippet1 { top: 15%; right: 8%; animation: floatRotate3D 8s ease-in-out infinite, glowPulseSnippet 4s ease-in-out infinite; }
.snippet2 { top: 55%; left: 5%; animation: floatRotate3D 10s ease-in-out infinite reverse, glowPulseSnippet 5s ease-in-out infinite 1s; }
.snippet3 { bottom: 20%; right: 12%; animation: floatRotate3D 12s ease-in-out infinite, glowPulseSnippet 6s ease-in-out infinite 2s; }

@keyframes floatRotate3D {
    0%, 100% { transform: translateY(0) rotateY(0deg); }
    50% { transform: translateY(-15px) rotateY(180deg); }
}

@keyframes glowPulseSnippet {
    0%, 100% { box-shadow: 0 10px 40px hsl(0 0% 0% / 0.4), 0 0 20px hsl(var(--primary) / 0.2); }
    50% { box-shadow: 0 10px 40px hsl(0 0% 0% / 0.5), 0 0 40px hsl(var(--primary) / 0.35); }
}

.heroWaves { position: absolute; bottom: 0; left: 0; width: 100%; height: 100px; animation: wave 10s linear infinite; opacity: 0.3; }

@keyframes wave {
    0% { transform: translateX(0) translateY(0); }
    50% { transform: translateX(-20px) translateY(-10px); }
    100% { transform: translateX(0) translateY(0); }
}

.heroContent {
    animation: fadeInUp 0.8s ease-out;
    position: relative; z-index: 100 !important;
    max-width: 1200px; margin: 0 auto;
    pointer-events: auto !important;
}

.heroTitle {
    font-size: 4.5rem; font-weight: 800; margin-bottom: 1.5rem;
    line-height: 1.1; letter-spacing: -1px;
    background: linear-gradient(90deg, #ffffff 0%, hsl(var(--primary)) 50%, #38bdf8 100%);
    background-clip: text; -webkit-background-clip: text;
    color: transparent; -webkit-text-fill-color: transparent;
    filter: drop-shadow(0 4px 16px rgba(0, 0, 0, 0.5));
    animation: slideDown 0.8s cubic-bezier(0.34, 1.56, 0.64, 1), windShimmer 6s ease-in-out infinite;
    background-size: 200% 100%;
}

@keyframes windShimmer { 0%, 100% { background-position: 0 50%; } 50% { background-position: 100% 50%; } }
@keyframes slideDown { from { opacity: 0; transform: translateY(-50px); } to { opacity: 1; transform: translateY(0); } }

.heroSubtitle {
    font-size: 1.75rem; margin-bottom: 1rem;
    max-width: 700px; margin-left: auto; margin-right: auto;
    opacity: 0.95; font-weight: 300; letter-spacing: 0.5px;
    background: linear-gradient(90deg, hsl(var(--primary)) 0%, #38bdf8 100%);
    background-clip: text; -webkit-background-clip: text;
    color: transparent; -webkit-text-fill-color: transparent;
    filter: drop-shadow(0 0 16px hsl(var(--primary) / 0.25));
    animation: fadeInUp 0.8s ease-out 0.2s both;
    line-height: 1.4;
}

.heroDescription {
    font-size: 1.15rem; margin-bottom: 2.5rem;
    max-width: 750px; margin-left: auto; margin-right: auto;
    opacity: 0.75; font-weight: 300; letter-spacing: 0.3px;
    color: hsl(0 0% 75%);
    animation: fadeInUp 0.8s ease-out 0.4s both;
    line-height: 1.6;
}

.heroActions {
    display: flex; gap: 1.5rem; justify-content: center; flex-wrap: wrap;
    animation: fadeInUp 0.8s ease-out 0.6s both;
    position: relative; z-index: 100 !important;
    pointer-events: auto !important;
}

@keyframes fadeInUp {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
}

.btnPrimary {
    background: linear-gradient(135deg, hsl(var(--primary)) 0%, hsl(var(--primary) / 0.85) 100%) !important;
    color: hsl(0 0% 98%) !important;
    font-size: 1.05rem !important; font-weight: 700;
    padding: 16px 40px !important; border-radius: 10px !important;
    border: none !important;
    box-shadow: 0 8px 24px hsl(var(--primary) / 0.35), 0 4px 12px hsl(var(--primary) / 0.2);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative; overflow: hidden;
    display: inline-flex !important; align-items: center;
}

.btnPrimary:hover {
    transform: translateY(-2px) translateX(4px) scale(1.02);
    box-shadow: 0 12px 32px hsl(var(--primary) / 0.5), 0 0 24px hsl(var(--primary) / 0.4);
}

.btnSecondary {
    background: hsl(0 0% 100% / 0.08) !important;
    color: hsl(0 0% 98%) !important;
    font-size: 1.05rem !important; font-weight: 600;
    padding: 16px 40px !important; border-radius: 10px !important;
    border: 1px solid hsl(var(--primary) / 0.35) !important;
    backdrop-filter: blur(12px);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative; overflow: hidden;
    display: inline-flex !important; align-items: center;
}

.btnSecondary:hover {
    background: hsl(var(--primary) / 0.12) !important;
    border-color: hsl(var(--primary) / 0.6) !important;
    transform: translateY(-2px) translateX(4px) scale(1.02);
}

@media (max-width: 900px) {
    .hero { padding: 7rem 1rem 2rem; }
    .heroContent { padding: 0 1rem; }
    .heroTitle { font-size: 2.5rem; }
    .heroSubtitle { font-size: 1.25rem; }
    .heroDescription { font-size: 1rem; }
    .heroCodeSnippets { display: none; }
    .heroGlow { width: 640px; height: 640px; filter: blur(120px); }
    .heroActions { flex-direction: column; gap: 0.75rem; padding: 0; }
    .btnPrimary, .btnSecondary { width: 100% !important; justify-content: center !important; }
}
</style>
