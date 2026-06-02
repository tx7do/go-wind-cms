import React from 'react';
import {Button} from '@/components/ui/button';
import {useTranslations} from 'next-intl';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

import styles from './home.module.css';

export default function HeroSection() {
    const t = useTranslations('page.home');
    const brandT = useTranslations('authentication.login');
    const router = useI18nRouter();

    return (
        <section className={styles.hero}>
            {/* Animated Background Elements */}
            <div className={styles.heroBgWrapper}>
                {/* 霓虹電能綠極光呼吸燈 */}
                <div className={styles.heroGlow}></div>
                <div className={styles.heroGradientBg}></div>
                {/* Grid Background */}
                <div className={styles.heroGridBg}></div>
                {/* Animated Shapes */}
                <div className={styles.heroAnimatedShapes}>
                    <div className={`${styles.shape} ${styles.shape1}`}></div>
                    <div className={`${styles.shape} ${styles.shape2}`}></div>
                    <div className={`${styles.shape} ${styles.shape3}`}></div>
                </div>
                {/* Code Snippet Decorations */}
                <div className={styles.heroCodeSnippets}>
                    <div className={`${styles.codeSnippet} ${styles.snippet1}`}>
                        <div className={styles.codeLine}>
                            <span className={styles.codeKeyword}>func</span>{' '}
                            <span className={styles.codeFunction}>GetPost</span>() {'{'}
                        </div>
                        <div className={styles.codeLine}>
                            <span className={styles.codeKeyword}>return</span> post
                        </div>
                        <div className={styles.codeLine}>{'}'}</div>
                    </div>
                    <div className={`${styles.codeSnippet} ${styles.snippet2}`}>
                        <div className={styles.codeLine}>
                            <span className={styles.codeTag}>&lt;template&gt;</span>
                        </div>
                        <div className={styles.codeLine}>
                            <span className={styles.codeTag}>&lt;div</span>{' '}
                            <span className={styles.codeAttr}>class</span>=
                            <span className={styles.codeString}>&quot;content-hub&quot;</span>
                            <span className={styles.codeTag}>&gt;</span>
                        </div>
                        <div className={styles.codeLine}>
                            <span className={styles.codeTag}>&lt;/div&gt;</span>
                        </div>
                    </div>
                    <div className={`${styles.codeSnippet} ${styles.snippet3}`}>
                        <div className={styles.codeLine}>
                            <span className={styles.codeComment}>{'// Content Hub API'}</span>
                        </div>
                        <div className={styles.codeLine}>
                            <span className={styles.codeKeyword}>const</span> api ={' '}
                            <span className={styles.codeString}>&apos;v1&apos;</span>
                        </div>
                    </div>
                </div>
                {/* SVG Waves */}
                <svg className={styles.heroWaves} viewBox="0 0 1440 320" preserveAspectRatio="none">
                    <path fill="hsl(var(--primary) / 0.08)"
                          d="M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,122.7C672,117,768,139,864,138.7C960,139,1056,117,1152,101.3C1248,85,1344,75,1392,69.3L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"></path>
                </svg>
            </div>

            <div className={styles.heroContent}>
                <h1 className={styles.heroTitle}>{brandT('brand_title')}</h1>
                <p className={styles.heroSubtitle}>{brandT('brand_subtitle')}</p>
                <p className={styles.heroDescription}>{t('hero_description')}</p>
                <div className={styles.heroActions}>
                    <Button size="lg" className={styles.btnPrimary}
                            onClick={() => router.push('/post')}>
                        {t('browse_posts')}
                    </Button>
                    <Button variant="outline" size="lg" className={styles.btnSecondary}
                            onClick={() => router.push('/about')}>
                        {t('learn_more')}
                    </Button>
                </div>
            </div>
        </section>
    );
}
