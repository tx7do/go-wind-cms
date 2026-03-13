'use client';

import {useTranslations} from 'next-intl';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './contact.module.css';

export default function ContactPage() {
    const t = useTranslations('page.legal.contact');

    return (
        <div className={styles.infoPage}>
            {/* Hero Section */}
            <section className={styles.hero}>
                <div className={styles.heroContent}>
                    <h1 className={styles.heroTitle}>{t('title')}</h1>
                    <p className={styles.heroSubtitle}>{t('description')}</p>
                </div>
            </section>

            {/* Content Section */}
            <div className={styles.infoCard}>
                <ul className={styles.list}>
                    <li>{t('item_1')}</li>
                    <li>{t('item_2')}</li>
                    <li>{t('item_3')}</li>
                </ul>
            </div>
        </div>
    );
}
