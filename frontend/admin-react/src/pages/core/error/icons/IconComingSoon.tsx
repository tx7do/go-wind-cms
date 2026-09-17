/**
 * Coming Soon 页面图标组件
 */
export const IconComingSoon = () => {
    return (
        <svg
            height="600"
            viewBox="0 0 600 600"
            width="600"
            xmlns="http://www.w3.org/2000/svg"
        >
            {/* Background circle */}
            <circle cx="300" cy="300" fill="var(--fb-mist)" r="280" />

            {/* Rocket body */}
            <path
                d="M300 150 C300 150 250 200 250 280 L250 350 C250 380 270 400 300 400 C330 400 350 380 350 350 L350 280 C350 200 300 150 300 150 Z"
                fill="var(--fb-navy)"
            />

            {/* Rocket window */}
            <circle cx="300" cy="280" fill="var(--fb-paper)" r="25" />
            <circle cx="300" cy="280" fill="var(--fb-primary)" r="15" />

            {/* Rocket fins */}
            <path
                d="M250 320 L220 380 L250 350 Z"
                fill="#ff6b6b"
            />
            <path
                d="M350 320 L380 380 L350 350 Z"
                fill="#ff6b6b"
            />

            {/* Flame */}
            <path
                d="M280 400 L300 450 L320 400 Z"
                fill="#ff9800"
            />
            <path
                d="M285 400 L300 435 L315 400 Z"
                fill="#ffeb3b"
            />

            {/* Stars */}
            <circle cx="150" cy="200" fill="#ffd700" r="4" />
            <circle cx="450" cy="180" fill="#ffd700" r="4" />
            <circle cx="180" cy="450" fill="#ffd700" r="4" />
            <circle cx="420" cy="420" fill="#ffd700" r="4" />
            <circle cx="120" cy="350" fill="#ffd700" r="3" />
            <circle cx="480" cy="300" fill="#ffd700" r="3" />

            {/* Text */}
            <text
                x="300"
                y="520"
                fontSize="40"
                fontWeight="bold"
                fill="var(--fb-navy)"
                textAnchor="middle"
                fontFamily="Arial, sans-serif"
            >
                Coming Soon
            </text>
        </svg>
    );
};
