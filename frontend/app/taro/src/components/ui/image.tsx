import React, {useState} from 'react';
import {Image as TaroImage} from '@tarojs/components';
import placeholderImage from '@/assets/images/placeholder.png';

export interface ImageProps {
    src?: string;
    fallbackSrc?: string;
    className?: string;
    mode?: 'scaleToFill' | 'aspectFit' | 'aspectFill' | 'widthFix' | 'heightFix' | 'top' | 'bottom' | 'center' | 'left' | 'right' | 'top left' | 'top right' | 'bottom left' | 'bottom right';
    style?: React.CSSProperties;
    onError?: (e: any) => void;
    [key: string]: any;
}

const Image: React.FC<ImageProps> = ({fallbackSrc = placeholderImage, onError, src, ...rest}) => {
    const [hasError, setHasError] = useState(false);

    const handleError = (e: any) => {
        if (!hasError) {
            setHasError(true);
        }
        onError?.(e);
    };

    return (
        <TaroImage
          src={hasError ? fallbackSrc : src}
          onError={handleError}
          mode='aspectFill'
          {...rest}
        />
    );
};

export default Image;
