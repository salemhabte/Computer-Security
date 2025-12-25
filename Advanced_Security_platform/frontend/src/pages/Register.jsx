import React, { useState } from 'react';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Card, CardHeader, CardBody } from '../components/ui/Card';
import { Link, useNavigate } from 'react-router-dom';
import { ShieldCheck, UserPlus } from 'lucide-react';
import api from '../lib/axios';

export default function Register() {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        username: '',
        phone: '',
        department: 'Sales', // Default, should be select
        role: 'USER'
    });
    const [captchaVerified, setCaptchaVerified] = useState(false);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        if (!captchaVerified) {
            setError('Please complete the CAPTCHA verification');
            return;
        }

        setLoading(true);
        try {
            // payload structure matches what UserUsecase.HandleRegistration expects
            // Actually backend expects: email, password, username, phone, department, role etc.
            // And captcha_token
            await api.post('/registration', {
                ...formData,
                captcha_token: 'dummy_token_verified'
            });

            // Navigate to verification page with email pre-filled (via state or query param)
            navigate('/verify-otp', { state: { email: formData.email } });

        } catch (err) {
            console.error(err);
            let errorMessage = 'Registration failed. Please try again.';
            if (err.response) {
                // Backend error
                errorMessage = err.response.data.error || errorMessage;
            } else if (err.request) {
                // Network error (no response)
                errorMessage = 'Network Error. Is the backend server running?';
            }
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4 py-8">
            <Card className="w-full max-w-md">
                <CardHeader className="text-center">
                    <div className="mx-auto w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mb-4">
                        <UserPlus className="w-6 h-6 text-green-600" />
                    </div>
                    <h1 className="text-2xl font-bold text-gray-900">Create Account</h1>
                    <p className="text-gray-500 mt-2">Join the secure platform</p>
                </CardHeader>
                <CardBody>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <Input
                            label="Username"
                            required
                            value={formData.username}
                            onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                            placeholder="johndoe"
                        />
                        <Input
                            label="Email"
                            type="email"
                            required
                            value={formData.email}
                            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                            placeholder="name@company.com"
                        />
                        <Input
                            label="Phone"
                            required
                            value={formData.phone}
                            onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                            placeholder="+1234567890"
                        />
                        <div className="w-full">
                            <label className="block text-sm font-medium text-gray-700 mb-1">Department</label>
                            <select
                                className="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                                value={formData.department}
                                onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                            >
                                <option value="Sales">Sales</option>
                                <option value="Marketing">Marketing</option>
                                <option value="Engineering">Engineering</option>
                                <option value="Finance">Finance</option>
                                <option value="HR">HR</option>
                            </select>
                        </div>

                        <Input
                            label="Password"
                            type="password"
                            required
                            value={formData.password}
                            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                            placeholder="Create a strong password"
                            className="mt-1"
                        />
                        <p className="text-xs text-gray-400">Must contain uppercase, lowercase, number, and special char.</p>

                        {/* Simulated CAPTCHA */}
                        <div className="flex items-center p-4 border rounded-lg bg-gray-50 mt-4">
                            <input
                                type="checkbox"
                                id="captcha-reg"
                                checked={captchaVerified}
                                onChange={(e) => setCaptchaVerified(e.target.checked)}
                                className="w-5 h-5 text-blue-600 rounded focus:ring-blue-500 border-gray-300"
                            />
                            <label htmlFor="captcha-reg" className="ml-3 flex items-center text-sm text-gray-700 select-none cursor-pointer">
                                <ShieldCheck className="w-4 h-4 mr-2 text-green-600" />
                                I am human (CAPTCHA)
                            </label>
                        </div>

                        {error && (
                            <div className="text-red-500 text-sm text-center bg-red-50 p-2 rounded">
                                {error}
                            </div>
                        )}

                        <Button type="submit" className="w-full mt-6" disabled={loading}>
                            {loading ? 'Creating Account...' : 'Sign Up'}
                        </Button>

                        <div className="text-center text-sm text-gray-500 pt-2">
                            Already have an account?{' '}
                            <Link to="/login" className="text-blue-600 hover:text-blue-500 font-medium">
                                Sign in
                            </Link>
                        </div>
                    </form>
                </CardBody>
            </Card>
        </div>
    );
}
