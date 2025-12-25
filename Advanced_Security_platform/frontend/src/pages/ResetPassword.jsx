import React, { useState } from 'react';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Card, CardHeader, CardBody } from '../components/ui/Card';
import { Link, useNavigate } from 'react-router-dom';
import { LockKeyhole, CheckCircle } from 'lucide-react';
import api from '../lib/axios';

export default function ResetPassword() {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        token: '',
        newPassword: '',
        confirmPassword: ''
    });
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#~$%^&*()_+|<>?:{}]).{8,}$/;
        if (!passwordRegex.test(formData.newPassword)) {
            setError('Password must contain uppercase, lowercase, number, and special character.');
            return;
        }

        if (formData.newPassword !== formData.confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        setLoading(true);
        try {
            await api.put('/reset_password', {
                token: formData.token,
                new_password: formData.newPassword
            });
            setSuccess(true);
            setTimeout(() => navigate('/login'), 3000);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to reset password.');
        } finally {
            setLoading(false);
        }
    };

    if (success) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
                <Card className="w-full max-w-md text-center">
                    <CardBody className="py-8">
                        <div className="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
                            <CheckCircle className="w-8 h-8 text-green-600" />
                        </div>
                        <h2 className="text-2xl font-bold text-gray-900 mb-2">Password Reset Successful</h2>
                        <p className="text-gray-600 mb-6">Your password has been updated. Redirecting to login...</p>
                        <Button onClick={() => navigate('/login')} className="w-full">
                            Go to Login
                        </Button>
                    </CardBody>
                </Card>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
            <Card className="w-full max-w-md">
                <CardHeader className="text-center">
                    <div className="mx-auto w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mb-4">
                        <LockKeyhole className="w-6 h-6 text-blue-600" />
                    </div>
                    <h1 className="text-2xl font-bold text-gray-900">Reset Password</h1>
                    <p className="text-gray-500 mt-2">Enter your code and new password</p>
                </CardHeader>
                <CardBody>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <Input
                            label="Reset Code (OTP)"
                            required
                            value={formData.token}
                            onChange={(e) => setFormData({ ...formData, token: e.target.value })}
                            placeholder="Enter the code sent to your email"
                        />

                        <Input
                            label="New Password"
                            type="password"
                            required
                            value={formData.newPassword}
                            onChange={(e) => setFormData({ ...formData, newPassword: e.target.value })}
                            placeholder="Create new password"
                        />

                        <Input
                            label="Confirm Password"
                            type="password"
                            required
                            value={formData.confirmPassword}
                            onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
                            placeholder="Confirm new password"
                        />
                        <p className="text-xs text-gray-400">Must contain uppercase, lowercase, number, and special char.</p>

                        {error && (
                            <div className="bg-red-50 text-red-700 p-3 rounded-lg text-sm text-center">
                                {error}
                            </div>
                        )}

                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading ? 'Resetting...' : 'Reset Password'}
                        </Button>
                    </form>
                </CardBody>
            </Card>
        </div>
    );
}
