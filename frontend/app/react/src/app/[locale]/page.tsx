'use client';

import React, {useEffect} from "react";
import {XIcon} from '@/plugins/xicon';
import {useTranslations} from 'next-intl';

import '../globals.css';
import styles from './page.module.css';

import HeroSection from '../../components/home/HeroSection';
import FeaturedPostsSection from '../../components/home/FeaturedPostsSection';
import CategoryListSection from '../../components/home/CategoryListSection';
import PopularTagsSection from '../../components/home/PopularTagsSection';
import LatestPostsSection from '../../components/home/LatestPostsSection';
import FeaturesSection from '../../components/home/FeaturesSection';

export default function Home() {
    const t = useTranslations('page.home');

    // Intersection Observer for scroll reveal
    useEffect(() => {
        const observer = new window.IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        entry.target.classList.add(styles.isVisible);
                        observer.unobserve(entry.target);
                    }
                });
            },
            {threshold: 0.1, rootMargin: "0px 0px -100px 0px"}
        );
        const elements = document.querySelectorAll(
            `.${styles.scrollRevealItem}`
        );
        elements.forEach((el) => observer.observe(el));
        return () => observer.disconnect();
    }, []);

    // Scroll to categories
    const scrollToCategories = () => {
        const section = document.querySelector(".categories-section");
        if (section) {
            section.scrollIntoView({behavior: "smooth", block: "start"});
        }
    };

    return (
        <div className={styles.page}>
            <HeroSection/>
            <FeaturedPostsSection/>
            <CategoryListSection/>
            <PopularTagsSection/>
            <LatestPostsSection/>
            <FeaturesSection/>
        </div>
    );
}

