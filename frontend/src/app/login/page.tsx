"use client";

import React, { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import api from '@/lib/api';
import { BookOpen, Mail, Lock, User, Eye, EyeOff, GraduationCap } from 'lucide-react';

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
          <div className="relative inline-block mb-8 md:mb-10">
            <div className="w-16 h-16 md:w-20 md:h-20 bg-white/10 backdrop-blur-md rounded-2xl flex items-center justify-center border border-white/20 shadow-xl mx-auto">
              <GraduationCap className="w-8 h-8 md:w-10 md:h-10 text-indigo-200" />
            </div>
            <div className="absolute -inset-2 rounded-2xl animate-pulse-ring bg-indigo-400/20" />
          </div>
          <h1 className="text-3xl md:text-4xl font-extrabold mb-4 leading-tight tracking-wide">
            School Management<br />
            <span className="text-indigo-200">Portal</span>
          </h1>
          <p className="text-indigo-200 text-sm md:text-base leading-relaxed tracking-wide mb-8">
            Manage classes, tasks, attendance, and student records — all in one modern platform.
          </p>
          <div className="flex flex-wrap justify-center gap-3 md:gap-4">
            {['📚 Academic Tasks', '✅ Attendance', '👨‍🎓 Student Profiles', '📊 Grades'].map(f => (
              <span key={f} className="bg-white/10 backdrop-blur-sm text-white text-[10px] md:text-xs font-medium px-3 py-1.5 md:px-4 md:py-2 rounded-full border border-white/20">
                {f}
              </span>
            ))}
          </div>
        </div>
      </div>

      <div className="w-full lg:w-1/2 flex items-center justify-center p-6 sm:p-8 lg:p-12 bg-slate-50">
        <div className="w-full max-w-md animate-fade-in-up">
          <div className="flex lg:hidden justify-center mb-8">
            <div className="w-12 h-12 bg-indigo-600 rounded-2xl flex items-center justify-center shadow-lg">
              <BookOpen size={24} className="text-white" />
            </div>
          </div>
          <div className="mb-8">
            <h2 className="text-2xl md:text-3xl font-extrabold text-slate-900 tracking-tight text-center lg:text-left">Welcome back</h2>
            <p className="text-slate-500 mt-2 md:mt-4 text-sm text-center lg:text-left">Sign in to your account to continue</p>
          </div>
          <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 md:p-8 lg:p-10">
            {error && (
              <div className="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-xl text-xs md:text-sm flex items-center gap-2 transition-all">
                <span className="text-red-500 shrink-0">⚠</span> 
                <span>{error}</span>
              </div>
            )}
            <form onSubmit={handleLogin} className="space-y-5">
              <div>
                <label className="block text-xs md:text-sm font-semibold text-slate-700 mb-1.5 md:mb-2">Username</label>
                <div className="relative">
                  <User size={16} className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input
                    type="text"
                    value={username}
                    onChange={e => setUsername(e.target.value)}
                    placeholder="johndoe"
                    required
                    className="w-full pl-11 pr-4 py-3 md:py-3.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-900 placeholder-slate-400 shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 hover:border-slate-400 transition-all"
                  />
                </div>
              </div>
              <div>
                <label className="block text-xs md:text-sm font-semibold text-slate-700 mb-1.5 md:mb-2">Email Address</label>
                <div className="relative">
                  <Mail size={16} className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input
                    type="email"
                    value={email}
                    onChange={e => setEmail(e.target.value)}
                    placeholder="you@example.com"
                    required
                    className="w-full pl-11 pr-4 py-3 md:py-3.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-900 placeholder-slate-400 shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 hover:border-slate-400 transition-all"
                  />
                </div>
              </div>
              <div>
                <div className="flex items-center justify-between mb-1.5 md:mb-2">
                  <label className="block text-xs md:text-sm font-semibold text-slate-700">Password</label>
                </div>
                <div className="relative">
                  <Lock size={16} className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                    placeholder="••••••••"
                    required
                    className="w-full pl-11 pr-12 py-3 md:py-3.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-900 placeholder-slate-400 shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 hover:border-slate-400 transition-all"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors"
                  >
                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
              </div>
              <div className="pt-2 text-red-500">
                 <button
                   type="submit"
                   disabled={isLoading}
                   className="w-full py-3 md:py-4 px-6 rounded-xl text-white font-semibold text-sm transition-all shadow-md hover:shadow-lg active:scale-[0.98] disabled:opacity-60 disabled:cursor-not-allowed"
                   style={{ background: isLoading ? '#818cf8' : 'linear-gradient(135deg, #4f46e5, #7c3aed)' }}
                 >
                   {isLoading ? (
                     <span className="flex items-center justify-center gap-2">
                       <svg className="animate-spin w-4 h-4 md:w-5 md:h-5" viewBox="0 0 24 24" fill="none">
                         <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                         <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V8H4z" />
                       </svg>
                       Signing in...
                     </span>
                   ) : 'Sign In →'}
                 </button>
              </div>
            </form>
          </div>
          <p className="text-center mt-8 md:mt-10 text-xs md:text-sm text-slate-500">
            Don&apos;t have an account?{' '}
            <Link href="/register" className="font-semibold text-indigo-600 hover:text-indigo-500 hover:underline transition-colors">
              Create one here
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
