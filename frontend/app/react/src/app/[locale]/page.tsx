'use client';

import React, {useEffect} from "react";
import {XIcon} from '@/plugins/xicon';
import {useTranslations} from 'next-intl';

import HeroSection from '../../components/home/HeroSection';
import FeaturedPostsSection from '../../components/home/FeaturedPostsSection';
import CategoryListSection from '../../components/home/CategoryListSection';
import PopularTagsSection from '../../components/home/PopularTagsSection';
import LatestPostsSection from '../../components/home/LatestPostsSection';
import FeaturesSection from '../../components/home/FeaturesSection';

export default function Home() {
    const t = useTranslations('page.home');

    // Intersection Observer for scroll reveal (supports dynamically added elements)
    useEffect(() => {
        const intersectionObserver = new window.IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        entry.target.classList.add('is-visible');
                        intersectionObserver.unobserve(entry.target);
                    }
                });
            },
            {threshold: 0.1, rootMargin: "0px 0px -100px 0px"}
        );

        // Observe existing elements
        const observeElements = (root: Element | Document = document) => {
            root.querySelectorAll('.scroll-reveal-item:not(.is-visible)').forEach((el) => {
                intersectionObserver.observe(el);
            });
        };
        observeElements();

        // Watch for dynamically added elements (e.g. after async data loads)
        const mutationObserver = new MutationObserver((mutations) => {
            for (const mutation of mutations) {
                mutation.addedNodes.forEach((node) => {
                    if (node.nodeType !== Node.ELEMENT_NODE) return;
                    const el = node as Element;
                    if (el.classList?.contains('scroll-reveal-item') && !el.classList.contains('is-visible')) {
                        intersectionObserver.observe(el);
                    }
                    observeElements(el);
                });
            }
        });
        mutationObserver.observe(document.body, {childList: true, subtree: true});

        return () => {
            intersectionObserver.disconnect();
            mutationObserver.disconnect();
        };
    }, []);

    // Scroll to categories
    const scrollToCategories = () => {
        const section = document.querySelector(".categories-section");
        if (section) {
            section.scrollIntoView({behavior: "smooth", block: "start"});
        }
    };

    return (
        <div className="w-full">
            <HeroSection/>
            <FeaturedPostsSection/>
            <CategoryListSection/>
            <PopularTagsSection/>
            <LatestPostsSection/>
            <FeaturesSection/>
        </div>
    );
}

