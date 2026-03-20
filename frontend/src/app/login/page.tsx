"use client";

import React, { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import api from '@/lib/api';
import { Mail, Lock, User, Eye, EyeOff, GraduationCap, BookOpen } from 'lucide-react';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);
    try {
      const response = await api.post('/login', { username, email, password });
      const payload = response.data?.data;
      if (payload && payload.token) {
        login(payload, payload.token);
      } else {
        setError('Login failed. Please check your credentials.');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || 'Invalid credentials. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex text-slate-900">
      <div className="hidden lg:flex lg:w-1/2 relative items-center justify-center p-8 lg:p-12 overflow-hidden" style={{ background: 'linear-gradient(135deg, #312e81 0%, #4f46e5 50%, #7c3aed 100%)' }}>
        <div className="absolute top-[-80px] left-[-80px] w-72 h-72 rounded-full opacity-20 animate-float" style={{ background: 'radial-gradient(circle, #818cf8, transparent)' }} />
        <div className="absolute bottom-[-60px] right-[-60px] w-96 h-96 rounded-full opacity-20 animate-float-delayed" style={{ background: 'radial-gradient(circle, #c084fc, transparent)' }} />
        <div className="absolute top-1/2 left-1/4 w-48 h-48 rounded-full opacity-10 animate-float" style={{ background: 'radial-gradient(circle, #38bdf8, transparent)', animationDelay: '3s' }} />
        <div className="absolute inset-0 opacity-5" style={{ backgroundImage: 'radial-gradient(circle, white 1px, transparent 1px)', backgroundSize: '30px 30px' }} />
        <div className="relative z-10 text-center text-white max-w-sm">
          <div className="relative inline-block mb-10">
            <div className="w-20 h-20 bg-white/10 backdrop-blur-md rounded-2xl flex items-center justify-center border border-white/20 shadow-xl mx-auto">
              <GraduationCap className="w-10 h-10 text-indigo-200" />
            </div>
            <div className="absolute -inset-2 rounded-2xl animate-pulse-ring bg-indigo-400/20" />
          </div>
          <h1 className="text-4xl font-extrabold mb-4 leading-tight tracking-wide">
            School Management<br />
            <span className="text-indigo-200">Portal</span>
          </h1>
          <p className="text-indigo-200 text-base leading-relaxed tracking-wide mb-8 font-medium opacity-90">
            Manage classes, tasks, attendance, and student records — all in one modern platform.
          </p>
          <div className="flex flex-wrap justify-center gap-4">
            {['📚 Academic Tasks', '✅ Attendance', '👨‍🎓 Student Profiles', '📊 Grades'].map(f => (
              <span key={f} className="bg-white/10 backdrop-blur-sm text-white text-xs font-bold px-4 py-2 rounded-full border border-white/20">
                {f}
              </span>
            ))}
          </div>
        </div>
      </div>

      <div className="w-full lg:w-1/2 flex items-center justify-center p-6 sm:p-8 lg:p-12 bg-slate-50">
        <div className="w-full max-w-md animate-fade-in-up">
          <div className="flex lg:hidden justify-center mb-10">
            <div className="w-16 h-16 bg-indigo-600 rounded-2xl flex items-center justify-center shadow-lg transform rotate-3 hover:rotate-0 transition-transform">
              <BookOpen size={32} className="text-white" />
            </div>
          </div>
          <div className="mb-10 space-y-2">
            <h2 className="text-3xl font-black text-slate-900 tracking-tight text-center lg:text-left">Welcome back</h2>
            <p className="text-slate-500 text-base text-center lg:text-left font-medium">Sign in to your account to continue</p>
          </div>
          <div className="bg-white rounded-3xl shadow-xl shadow-slate-200/50 border border-slate-200 p-8 md:p-10">
            {error && (
              <div className="mb-8 p-4 bg-red-50 border border-red-100 text-red-600 rounded-2xl text-sm font-semibold flex items-center gap-3 animate-head-shake">
                <span className="text-lg">⚠</span> 
                <span>{error}</span>
              </div>
            )}
            <form onSubmit={handleLogin} className="space-y-6">
              <Input
                label="Username"
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                placeholder="johndoe"
                required
                icon={<User size={18} />}
              />
              <Input
                label="Email Address"
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                placeholder="you@example.com"
                required
                icon={<Mail size={18} />}
              />
              <div className="relative">
                <Input
                  label="Password"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  placeholder="••••••••"
                  required
                  icon={<Lock size={18} />}
                  className="pr-12"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-4 top-[38px] text-slate-400 hover:text-indigo-600 transition-colors z-20"
                >
                  {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
              
              <div className="pt-4">
                <Button
                  type="submit"
                  isLoading={isLoading}
                  fullWidth
                  size="lg"
                  className="font-black tracking-wide"
                >
                  Sign In →
                </Button>
              </div>
            </form>
          </div>
          <p className="text-center mt-10 text-sm text-slate-500 font-medium">
            Don&apos;t have an account?{' '}
            <Link href="/register" className="font-black text-indigo-600 hover:text-indigo-500 hover:underline transition-all">
              Create one here
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
