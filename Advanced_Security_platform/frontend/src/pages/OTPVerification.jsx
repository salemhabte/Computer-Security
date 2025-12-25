import React, { useState, useEffect } from 'react';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Card, CardHeader, CardBody } from '../components/ui/Card';
import { useNavigate, useLocation } from 'react-router-dom';
import { CheckCircle } from 'lucide-react';
import api from '../lib/axios';

export default function OTPVerification() {
    const navigate = useNavigate();
    const location = useLocation();
    const [email, setEmail] = useState('');
    const [otp, setOtp] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState(false);

    useEffect(() => {
        if (location.state?.email) {
            setEmail(location.state.email);
        }
    }, [location.state]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            await api.post('/registration/verification', {
                email,
                otp
            });
            setSuccess(true);
            setTimeout(() => navigate('/login'), 2000);
        } catch (err) {
            setError(err.response?.data?.error || 'Verification failed. Invalid OTP or expired.');
        } finally {
            setLoading(false);
        }
    };

    if (success) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
                <Card className="w-full max-w-md text-center p-8">
                    <div className="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-6">
                        <CheckCircle className="w-10 h-10 text-green-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900 mb-2">Verified!</h2>
                    <p className="text-gray-500">Your account has been successfully verified.</p>
                    <p className="text-sm text-gray-400 mt-4">Redirecting to login...</p>
                </Card>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
            <Card className="w-full max-w-md">
                <CardHeader className="text-center">
                    <h1 className="text-2xl font-bold text-gray-900">Verify Email</h1>
                    <p className="text-gray-500 mt-2">Enter the code sent to {email}</p>
                </CardHeader>
                <CardBody>
                    <form onSubmit={handleSubmit} className="space-y-6">
                        <Input
                            label="OTP Code"
                            required
                            value={otp}
                            onChange={(e) => setOtp(e.target.value)}
                            placeholder="123456"
                            className="text-center tracking-[0.5em] text-xl font-mono"
                        />

                        {error && (
                            <div className="text-red-500 text-sm text-center bg-red-50 p-2 rounded">
                                {error}
                            </div>
                        )}

                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading ? 'Verifying...' : 'Verify Email'}
                        </Button>

                        <div className="text-center text-sm">
                            <button type="button" className="text-blue-600 hover:underline" onClick={() => navigate('/register')}>
                                Wrong email? Back to register
                            </button>
                        </div>
                    </form>
                </CardBody>
            </Card>
        </div>
    );
}
