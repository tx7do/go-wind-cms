import {useState, useEffect} from 'react';
import {useTranslations} from 'next-intl';
import styles from '../register.scss';

export default function PhoneRegisterPage() {
    const t = useTranslations('authentication');

    const [phone, setPhone] = useState('');
    const [verificationCode, setVerificationCode] = useState('');
    const [codeSent, setCodeSent] = useState(false);
    const [codeCountdown, setCodeCountdown] = useState(0);

    useEffect(() => {
        let codeTimer: number | null = null;

        if (codeSent && codeCountdown > 0) {
            codeTimer = window.setInterval(() => {
                setCodeCountdown((prev) => {
                    if (prev <= 1) {
                        setCodeSent(false);
                        return 0;
                    }
                    return prev - 1;
                });
            }, 1000);
        }

        return () => {
            if (codeTimer !== null) {
                clearInterval(codeTimer);
            }
        };
    }, [codeSent, codeCountdown]);

    /**
     * 发送验证码
     */
    const handleButtonSendVerifyCode = () => {
        if (!phone) {
            return;
        }

        setCodeSent(true);
        setCodeCountdown(60);
    };

    /**
     * 登录或者注册
     */
    const handleButtonRegisterOrLogin = () => {
        if (!phone || !verificationCode) {
            return;
        }
        console.log('手机号注册/登录:', {phone, verificationCode});
    };

    return (
        <div className={styles['register-form']}>
            {/* Phone Number Group */}
            <div className={styles['form-group']}>
                <label htmlFor="register-phone-number">
                    {t('register.phone')}
                </label>
                <input
                    id="register-phone-number"
                    type="text"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    placeholder={t('register.input_phone')}
                    autoComplete="tel"
                    className={styles['input-field']}
                />
            </div>

            {/* Verification Code Group */}
            <div className={styles['form-group']}>
                <label htmlFor="register-phone-code">
                    {t('register.code')}
                </label>
                <div className={styles['code-input-row']}>
                    <input
                        id="register-phone-code"
                        type="text"
                        value={verificationCode}
                        onChange={(e) => setVerificationCode(e.target.value)}
                        placeholder={t('register.input_code')}
                        maxLength={6}
                        autoComplete="one-time-code"
                        className={styles['input-field']}
                    />
                    <button
                        type="button"
                        disabled={codeSent}
                        className={`${styles['send-code-btn']} ${codeSent ? styles.disabled : ''}`}
                        onClick={handleButtonSendVerifyCode}
                    >
                        {codeSent ? `${codeCountdown}s` : t('register.send_code')}
                    </button>
                </div>
            </div>

            {/* Register Button */}
            <button
                type="button"
                className={styles['register-button']}
                onClick={handleButtonRegisterOrLogin}
            >
                {t('register.register_or_login')}
            </button>
        </div>
    );
}
