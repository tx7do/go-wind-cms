import * as React from 'react';
import {View, ScrollView} from '@tarojs/components';
import {cn} from '@/lib/utils';
import {Button} from '@/components/ui/button';
import {XIcon} from '@/plugins/xicon';

interface CarouselProps {
    orientation?: 'horizontal' | 'vertical';
    className?: string;
    children: React.ReactNode;
    opts?: any;
    setApi?: (api: any) => void;
    plugins?: any;
}

interface CarouselContextProps {
    orientation: 'horizontal' | 'vertical';
    scrollPrev: () => void;
    scrollNext: () => void;
    canScrollPrev: boolean;
    canScrollNext: boolean;
    currentIndex: number;
    totalItems: number;
}

const CarouselContext = React.createContext<CarouselContextProps | null>(null);

function useCarousel() {
    const context = React.useContext(CarouselContext);
    if (!context) {
        throw new Error('useCarousel must be used within a <Carousel />');
    }
    return context;
}

function Carousel({orientation = 'horizontal', className, children, ...props}: CarouselProps) {
    const [currentIndex, setCurrentIndex] = React.useState(0);
    const childArray = React.Children.toArray(children);
    const totalItems = childArray.length;

    const scrollPrev = () => setCurrentIndex(i => Math.max(0, i - 1));
    const scrollNext = () => setCurrentIndex(i => Math.min(totalItems - 1, i + 1));

    return (
        <CarouselContext.Provider value={{
            orientation,
            scrollPrev,
            scrollNext,
            canScrollPrev: currentIndex > 0,
            canScrollNext: currentIndex < totalItems - 1,
            currentIndex,
            totalItems,
        }}
        >
            <View className={cn('relative', className)} {...props}>
                {children}
            </View>
        </CarouselContext.Provider>
    );
}
Carousel.displayName = 'Carousel';

function CarouselContent({className, children, ...props}: {className?: string; children: React.ReactNode; [key: string]: any}) {
    const {orientation} = useCarousel();
    return (
        <ScrollView
          scrollX={orientation === 'horizontal'}
          scrollY={orientation === 'vertical'}
          className={cn('overflow-hidden', className)}
          {...props}
        >
            <View className={cn('flex', orientation === 'horizontal' ? 'flex-row' : 'flex-col')}>
                {children}
            </View>
        </ScrollView>
    );
}
CarouselContent.displayName = 'CarouselContent';

function CarouselItem({className, children, ...props}: {className?: string; children: React.ReactNode; [key: string]: any}) {
    const {orientation} = useCarousel();
    return (
        <View
          className={cn(
                'min-w-0 shrink-0 grow-0 basis-full',
                orientation === 'horizontal' ? 'pl-4' : 'pt-4',
                className,
            )}
          {...props}
        >
            {children}
        </View>
    );
}
CarouselItem.displayName = 'CarouselItem';

function CarouselPrevious({className, variant, size}: {className?: string; variant?: any; size?: any}) {
    const {orientation, scrollPrev, canScrollPrev} = useCarousel();
    return (
        <Button
          variant={variant || 'outline'}
          size={size || 'icon'}
          className={cn(
                'absolute h-8 w-8 rounded-full',
                orientation === 'horizontal'
                    ? '-left-12 top-1/2 -translate-y-1/2'
                    : '-top-12 left-1/2 -translate-x-1/2 rotate-90',
                className,
            )}
          onClick={scrollPrev}
          disabled={!canScrollPrev}
        >
            <XIcon name='carbon:arrow-left' size={16} />
        </Button>
    );
}
CarouselPrevious.displayName = 'CarouselPrevious';

function CarouselNext({className, variant, size}: {className?: string; variant?: any; size?: any}) {
    const {orientation, scrollNext, canScrollNext} = useCarousel();
    return (
        <Button
          variant={variant || 'outline'}
          size={size || 'icon'}
          className={cn(
                'absolute h-8 w-8 rounded-full',
                orientation === 'horizontal'
                    ? '-right-12 top-1/2 -translate-y-1/2'
                    : '-bottom-12 left-1/2 -translate-x-1/2 rotate-90',
                className,
            )}
          onClick={scrollNext}
          disabled={!canScrollNext}
        >
            <XIcon name='carbon:arrow-right' size={16} />
        </Button>
    );
}
CarouselNext.displayName = 'CarouselNext';

export {
    Carousel,
    CarouselContent,
    CarouselItem,
    CarouselPrevious,
    CarouselNext,
};
