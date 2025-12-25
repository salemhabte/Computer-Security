import React, { useState } from 'react';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Card, CardHeader, CardBody } from '../components/ui/Card';
import { Link, useNavigate } from 'react-router-dom';
import { KeyRound, ArrowLeft } from 'lucide-react';
import api from '../lib/axios';

export default function ForgotPassword() {
    const navigate = useNavigate();
    const [email, setEmail] = useState('');
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setMessage('');
        setError('');

        try {
            await api.post('/forgot_password', { email });
            setMessage('Password reset instructions have been sent to your email.');
            // Optionally navigate to reset page after a delay
            setTimeout(() => navigate('/reset-password'), 2000);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to send reset link.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
            <Card className="w-full max-w-md">
                <CardHeader className="text-center">
                    <div className="mx-auto w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mb-4">
                        <KeyRound className="w-6 h-6 text-blue-600" />
                    </div>
                    <h1 className="text-2xl font-bold text-gray-900">Forgot Password?</h1>
                    <p className="text-gray-500 mt-2">Enter your email to receive a reset code</p>
                </CardHeader>
                <CardBody>
                    <form onSubmit={handleSubmit} className="space-y-6">
                        <Input
                            label="Email Address"
                            type="email"
                            required
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder="name@company.com"
                        />

                        {message && (
                            <div className="bg-green-50 text-green-700 p-3 rounded-lg text-sm text-center">
                                {message}
                            </div>
                        )}

                        {error && (
                            <div className="bg-red-50 text-red-700 p-3 rounded-lg text-sm text-center">
                                {error}
                            </div>
                        )}

                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading ? 'Sending Link...' : 'Send Reset Link'}
                        </Button>

                        <div className="text-center">
                            <Link to="/login" className="text-gray-600 hover:text-gray-900 text-sm flex items-center justify-center gap-2">
                                <ArrowLeft size={16} />
                                Back to Login
                            </Link>
                        </div>
                    </form>
                </CardBody>
            </Card>
        </div>
    );
}
