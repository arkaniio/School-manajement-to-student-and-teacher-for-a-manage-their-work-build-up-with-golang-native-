"use client";

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import api from '@/lib/api';
import { Mail, Lock, User, Eye, EyeOff, GraduationCap, BookOpen } from 'lucide-react';

export default function RegisterPage() {
  const [formData, setFormData] = useState({ username: '', email: '', password: '', role: 'siswa' });
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setIsLoading(true);
    try {
      await api.post('/register', formData);
      setSuccess('Registration successful! Redirecting to login...');
      setTimeout(() => router.push('/login'), 2000);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Registration failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex text-slate-900">
      {/* Left — Premium Branding Panel */}
      <div className="hidden lg:flex lg:w-1/2 relative items-center justify-center p-8 lg:p-12 overflow-hidden"
        style={{ background: 'linear-gradient(135deg, #064e3b 0%, #059669 50%, #0d9488 100%)' }}>

        <div className="absolute top-[-80px] right-[-60px] w-80 h-80 rounded-full opacity-20 animate-float"
          style={{ background: 'radial-gradient(circle, #34d399, transparent)' }} />
        <div className="absolute bottom-[-40px] left-[-40px] w-72 h-72 rounded-full opacity-20 animate-float-delayed"
          style={{ background: 'radial-gradient(circle, #2dd4bf, transparent)' }} />
        <div className="absolute inset-0 opacity-5"
          style={{ backgroundImage: 'radial-gradient(circle, white 1px, transparent 1px)', backgroundSize: '30px 30px' }} />

        <div className="relative z-10 text-center text-white max-w-sm">
          <div className="relative inline-block mb-8 md:mb-10">
            <div className="w-16 h-16 md:w-20 md:h-20 bg-white/10 backdrop-blur-md rounded-2xl flex items-center justify-center border border-white/20 shadow-xl mx-auto">
              <GraduationCap className="w-8 h-8 md:w-10 md:h-10 text-emerald-200" />
            </div>
            <div className="absolute -inset-2 rounded-2xl animate-pulse-ring bg-emerald-400/20" />
          </div>
          <h1 className="text-3xl md:text-4xl font-extrabold mb-4 leading-normal tracking-wide">
            Join Our<br /><span className="text-emerald-200">Learning Platform</span>
          </h1>
          <p className="text-emerald-100 text-sm md:text-base leading-relaxed tracking-wide mb-8">
            Create your student or teacher account and start your journey with a modern school management system.
          </p>
          <div className="flex flex-wrap justify-center gap-3 md:gap-4">
            {['🎓 Student Profiles', '📝 Task Management', '📅 Attendance', '🏆 Grading'].map(f => (
              <span key={f} className="bg-white/10 backdrop-blur-sm text-white text-[10px] md:text-xs font-medium px-3 py-1.5 md:px-4 md:py-2 rounded-full border border-white/20">
                {f}
              </span>
            ))}
          </div>
        </div>
      </div>

      {/* Right — Form Panel */}
      <div className="w-full lg:w-1/2 flex items-center justify-center p-6 sm:p-8 lg:p-12 bg-slate-50">
        <div className="w-full max-w-md animate-fade-in-up mt-8 mb-8 lg:m-0">
          <div className="flex lg:hidden justify-center mb-8">
            <div className="w-12 h-12 bg-emerald-600 rounded-2xl flex items-center justify-center shadow-lg">
              <BookOpen size={24} className="text-white" />
            </div>
          </div>

          <div className="mb-6">
            <h2 className="text-2xl md:text-3xl font-extrabold text-slate-900 tracking-tight text-center lg:text-left">Create an account</h2>
            <p className="text-slate-500 mt-2 text-sm text-center lg:text-left">Fill in your details to get started</p>
          </div>

          <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 md:p-8">
            {error && (
              <div className="mb-5 p-4 bg-red-50 border border-red-200 text-red-700 rounded-xl text-xs md:text-sm flex items-center gap-4 transition-all">
                <span className="shrink-0 text-red-500">⚠</span>
                <span>{error}</span>
              </div>
            )}
            {success && (
              <div className="mb-5 p-4 bg-emerald-50 border border-emerald-200 text-emerald-700 rounded-xl text-xs md:text-sm flex items-center gap-4 transition-all">
                <span className="shrink-0 text-emerald-500">✅</span>
                <span>{success}</span>
              </div>
            )}

            <form onSubmit={handleRegister} className="space-y-4 md:space-y-5">
              {/* Username */}
              <div>
                <label className="block text-xs md:text-sm font-semibold text-slate-700 mb-1.5 md:mb-2">Full Name / Username</label>
                <div className="relative">
                  <User size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input name="username" type="text" value={formData.username} onChange={handleChange}
                    placeholder="John Doe" required
                    className="w-full pl-10 pr-4 py-3 md:py-3.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-900 placeholder-slate-400 shadow-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 hover:border-slate-400 transition-all" />
                </div>
              </div>

              {/* Email */}
              <div>
                <label className="block text-xs md:text-sm font-semibold text-slate-700 mb-1.5 md:mb-2">Email Address</label>
                <div className="relative">
                  <Mail size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input name="email" type="email" value={formData.email} onChange={handleChange}
                    placeholder="you@example.com" required
                    className="w-full pl-10 pr-4 py-3 md:py-3.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-900 placeholder-slate-400 shadow-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 hover:border-slate-400 transition-all" />
                </div>
              </div>

              {/* Password */}
              <div>
                <label className="block text-xs md:text-sm font-semibold text-slate-700 mb-1.5 md:mb-2">Password</label>
                <div className="relative">
                  <Lock size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input name="password" type={showPassword ? 'text' : 'password'} value={formData.password} onChange={handleChange}
                    placeholder="••••••••" required minLength={6}
                    className="w-full pl-10 pr-10 py-3 md:py-3.5 bg-white border border-slate-300 rounded-xl text-sm text-slate-900 placeholder-slate-400 shadow-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 hover:border-slate-400 transition-all" />
                  <button type="button" onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors">
                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
              </div>

              {/* Role Selector */}
              <div>
                <label className="block text-xs md:text-sm font-semibold text-slate-700 mb-3 md:mb-4">Account Role</label>
                <div className="grid grid-cols-2 gap-3 md:gap-4">
                  {[
                    { value: 'siswa', label: 'Student', emoji: '👨‍🎓', desc: 'Register as a student' },
                    { value: 'guru', label: 'Teacher', emoji: '👨‍🏫', desc: 'Register as a teacher' },
                  ].map(role => (
                    <label key={role.value}
                      className={`relative flex flex-col items-center text-center p-3 md:p-4 rounded-xl border-2 cursor-pointer transition-all ${
                        formData.role === role.value
                          ? 'border-emerald-500 bg-emerald-50 shadow-sm md:shadow-md'
                          : 'border-slate-200 bg-white hover:border-slate-300 hover:bg-slate-50'
                      }`}>
                      <input type="radio" name="role" value={role.value}
                        checked={formData.role === role.value}
                        onChange={handleChange}
                        className="sr-only" />
                      <span className="text-xl md:text-2xl mb-1 md:mb-2">{role.emoji}</span>
                      <span className={`text-[11px] md:text-sm font-bold ${formData.role === role.value ? 'text-emerald-700' : 'text-slate-700'}`}>
                        {role.label}
                      </span>
                      <span className="hidden md:block text-[10px] md:text-xs text-slate-400 mt-2">{role.desc}</span>
                      {formData.role === role.value && (
                        <span className="absolute top-2 right-2 md:top-3 md:right-3 w-3 h-3 md:w-4 md:h-4 bg-emerald-500 rounded-full flex items-center justify-center text-white text-[8px] md:text-[10px]">✓</span>
                      )}
                    </label>
                  ))}
                </div>
              </div>

              {/* Submit */}
              <div className="pt-2">
                <button
                  type="submit"
                  disabled={isLoading}
                  className="w-full py-3 md:py-3.5 px-6 rounded-xl text-white font-semibold text-sm transition-all shadow-md hover:shadow-lg active:scale-[0.98] disabled:opacity-60"
                  style={{ background: isLoading ? '#6ee7b7' : 'linear-gradient(135deg, #059669, #0d9488)' }}
                >
                  {isLoading ? (
                    <span className="flex items-center justify-center gap-2">
                      <svg className="animate-spin w-4 h-4 md:w-5 md:h-5" viewBox="0 0 24 24" fill="none">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                      </svg>
                      Creating account...
                    </span>
                  ) : 'Create Account →'}
                </button>
              </div>
            </form>
          </div>

          <p className="text-center mt-6 text-xs md:text-sm text-slate-500">
            Already have an account?{' '}
            <Link href="/login" className="font-semibold text-emerald-600 hover:text-emerald-500 hover:underline transition-colors">
              Sign in instead
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
