/**
 * Offline 页面图标组件
 */
export const IconOffline = () => {
    return (
        <svg
            height="600"
            viewBox="0 0 600 600"
            width="600"
            xmlns="http://www.w3.org/2000/svg"
        >
            {/* Background circle */}
            <circle cx="300" cy="300" fill="var(--fb-mist)" r="280" />

            {/* WiFi signal - disconnected */}
            <path
                d="M200 250 Q300 150 400 250"
                fill="none"
                stroke="var(--fb-mist-2)"
                strokeLinecap="round"
                strokeWidth="20"
            />
            <path
                d="M220 280 Q300 200 380 280"
                fill="none"
                stroke="var(--fb-mist-2)"
                strokeLinecap="round"
                strokeWidth="20"
            />
            <path
                d="M240 310 Q300 250 360 310"
                fill="none"
                stroke="var(--fb-mist-2)"
                strokeLinecap="round"
                strokeWidth="20"
            />

            {/* WiFi dot */}
            <circle cx="300" cy="340" fill="var(--fb-mist-2)" r="15" />

            {/* Cross mark */}
            <line
                x1="270"
                x2="330"
                y1="310"
                y2="370"
                stroke="#ff6b6b"
                strokeLinecap="round"
                strokeWidth="8"
            />
            <line
                x1="330"
                x2="270"
                y1="310"
                y2="370"
                stroke="#ff6b6b"
                strokeLinecap="round"
                strokeWidth="8"
            />

            {/* Computer/Device */}
            <rect
                fill="var(--fb-navy)"
                height="120"
                rx="10"
                ry="10"
                width="180"
                x="210"
                y="400"
            />

            {/* Screen */}
            <rect
                fill="var(--fb-paper)"
                height="90"
                width="160"
                x="220"
                y="410"
            />

            {/* Stand */}
            <rect
                fill="var(--fb-navy)"
                height="30"
                width="40"
                x="280"
                y="520"
            />
            <rect
                fill="var(--fb-navy)"
                height="10"
                width="80"
                x="260"
                y="550"
            />

            {/* Error indicator on screen */}
            <text
                x="300"
                y="465"
                fontSize="50"
                fontWeight="bold"
                fill="#ff6b6b"
                textAnchor="middle"
                fontFamily="Arial, sans-serif"
            >
                ×
            </text>
        </svg>
    );
};
