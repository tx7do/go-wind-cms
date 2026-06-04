import React from 'react';
import {View, Text} from '@tarojs/components';
import {Button} from '@/components/ui/button';
import {useTranslations} from '@/lib/next-intl-compat';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {XIcon} from '@/plugins/xicon';

export default function HeroSection() {
    const t = useTranslations('page.home');
    const brandT = useTranslations('authentication.login');
    const router = useI18nRouter();

    return (
        <View className='relative min-h-[500px] flex items-center justify-center overflow-hidden bg-gradient-to-br from-primary/90 via-primary to-accent/80 text-primary-foreground'>
            {/* Animated Background */}
            <View className='absolute inset-0'>
                {/* Glow */}
                <View className='absolute top-1/4 left-1/4 h-64 w-64 rounded-full bg-accent/20 blur-3xl' />
                <View className='absolute bottom-1/4 right-1/4 h-48 w-48 rounded-full bg-primary/30 blur-3xl' />
                {/* Grid Background */}
                <View className='absolute inset-0 opacity-[0.03]'
                  style={{backgroundImage: 'linear-gradient(rgba(255,255,255,.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,.1) 1px, transparent 1px)', backgroundSize: '60px 60px'}}
                />
                {/* Floating Shapes */}
                <View className='absolute top-20 left-[15%] h-16 w-16 rounded-lg border border-primary-foreground/10 bg-primary-foreground/5 animate-float-slow' />
                <View className='absolute top-32 right-[20%] h-12 w-12 rounded-full border border-primary-foreground/10 bg-primary-foreground/5 animate-float-medium' />
                <View className='absolute bottom-24 left-[30%] h-20 w-20 rounded-xl border border-primary-foreground/10 bg-primary-foreground/5 animate-float-fast' />
                {/* Code Snippet Decorations */}
                <View className='absolute top-16 right-[8%] hidden rounded-lg border border-primary-foreground/10 bg-background/10 px-3 py-2 font-mono text-xs backdrop-blur-sm md:block'>
                    <View className='flex gap-1'>
                        <Text className='text-accent'>func</Text>
                        <Text className='text-primary-foreground'>GetPost</Text>() {'{'}
                    </View>
                    <View className='flex gap-1'>
                        <Text className='text-accent'>return</Text> post
                    </View>
                    {'}'}
                </View>
                <View className='absolute bottom-20 left-[6%] hidden rounded-lg border border-primary-foreground/10 bg-background/10 px-3 py-2 font-mono text-xs backdrop-blur-sm md:block'>
                    <View className='flex gap-1'>
                        <Text className='text-accent'>const</Text> api ={' '}
                        <Text className='text-green-300'>&apos;v1&apos;</Text>
                    </View>
                </View>
            </View>

            {/* Content */}
            <View className='relative z-10 flex flex-col items-center gap-6 px-4 py-20 text-center'>
                <View className='mb-2 flex h-16 w-16 items-center justify-center rounded-2xl bg-primary-foreground/15 backdrop-blur-sm'>
                    <XIcon name='carbon:blog' size={32} className='text-primary-foreground' />
                </View>
                <Text className='text-4xl font-bold tracking-tight text-primary-foreground md:text-5xl lg:text-6xl'>
                    {brandT('brand_title')}
                </Text>
                <Text className='max-w-2xl text-lg text-primary-foreground/90 md:text-xl'>
                    {brandT('brand_subtitle')}
                </Text>
                <Text className='max-w-3xl text-base text-primary-foreground/80 md:text-lg'>
                    {t('hero_description')}
                </Text>
                <View className='mt-4 flex flex-wrap items-center justify-center gap-4'>
                    <Button size='lg'
                      className='bg-primary-foreground text-primary hover:bg-primary-foreground/90 shadow-lg'
                      onClick={() => router.push('/post')}
                    >
                        {t('browse_posts')}
                    </Button>
                    <Button variant='outline' size='lg'
                      className='border-primary-foreground/30 text-primary-foreground hover:bg-primary-foreground/10'
                      onClick={() => router.push('/about')}
                    >
                        {t('learn_more')}
                    </Button>
                </View>
            </View>
        </View>
    );
}
