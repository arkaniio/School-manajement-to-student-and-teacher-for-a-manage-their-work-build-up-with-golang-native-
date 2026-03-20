"use client";

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import api from '@/lib/api';
import { Mail, Lock, User, Eye, EyeOff, GraduationCap, BookOpen, Check } from 'lucide-react';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';

export default function RegisterPage() {
  const [formData, setFormData] = useState({ username: '', email: '', password: '', role: 'siswa' });
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
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
          <div className="relative inline-block mb-10">
            <div className="w-20 h-20 bg-white/10 backdrop-blur-md rounded-2xl flex items-center justify-center border border-white/20 shadow-xl mx-auto">
              <GraduationCap className="w-10 h-10 text-emerald-200" />
            </div>
            <div className="absolute -inset-2 rounded-2xl animate-pulse-ring bg-emerald-400/20" />
          </div>
          <h1 className="text-4xl font-extrabold mb-4 leading-normal tracking-wide">
            Join Our<br /><span className="text-emerald-200">Learning Platform</span>
          </h1>
          <p className="text-emerald-100 text-base leading-relaxed tracking-wide mb-8 font-medium opacity-90">
            Create your student or teacher account and start your journey with a modern school management system.
          </p>
          <div className="flex flex-wrap justify-center gap-4">
            {['🎓 Student Profiles', '📝 Task Management', '📅 Attendance', '🏆 Grading'].map(f => (
              <span key={f} className="bg-white/10 backdrop-blur-sm text-white text-xs font-bold px-4 py-2 rounded-full border border-white/20">
                {f}
              </span>
            ))}
          </div>
        </div>
      </div>

      {/* Right — Form Panel */}
      <div className="w-full lg:w-1/2 flex items-center justify-center p-6 sm:p-8 lg:p-12 bg-slate-50">
        <div className="w-full max-w-md animate-fade-in-up mt-8 mb-8 lg:m-0">
          <div className="flex lg:hidden justify-center mb-10">
            <div className="w-16 h-16 bg-emerald-600 rounded-2xl flex items-center justify-center shadow-lg transform rotate-3 hover:rotate-0 transition-all">
              <BookOpen size={32} className="text-white" />
            </div>
          </div>

          <div className="mb-10 space-y-2">
            <h2 className="text-3xl font-black text-slate-900 tracking-tight text-center lg:text-left uppercase">Create account</h2>
            <p className="text-slate-500 text-base text-center lg:text-left font-medium">Fill in your details to get started</p>
          </div>

          <div className="bg-white rounded-3xl shadow-xl shadow-slate-200/50 border border-slate-200 p-8 md:p-10">
            {error && (
              <div className="mb-8 p-4 bg-red-50 border border-red-100 text-red-600 rounded-2xl text-sm font-semibold flex items-center gap-3 animate-head-shake">
                <span className="text-lg">⚠</span>
                <span>{error}</span>
              </div>
            )}
            {success && (
              <div className="mb-8 p-4 bg-emerald-50 border border-emerald-100 text-emerald-600 rounded-2xl text-sm font-semibold flex items-center gap-3 animate-fade-in">
                <span className="text-lg">✅</span>
                <span>{success}</span>
              </div>
            )}

            <form onSubmit={handleRegister} className="space-y-6">
              <Input
                label="Full Name / Username"
                name="username"
                type="text"
                value={formData.username}
                onChange={handleChange}
                placeholder="John Doe"
                required
                icon={<User size={18} />}
              />

              <Input
                label="Email Address"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                placeholder="you@example.com"
                required
                icon={<Mail size={18} />}
              />

              <div className="relative">
                <Input
                  label="Password"
                  name="password"
                  type={showPassword ? 'text' : 'password'}
                  value={formData.password}
                  onChange={handleChange}
                  placeholder="••••••••"
                  required
                  icon={<Lock size={18} />}
                  minLength={6}
                />
                <button type="button" onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-4 top-[38px] text-slate-400 hover:text-emerald-600 transition-colors z-20">
                  {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>

              <div>
                <label className="form-label block mb-4">Account Role</label>
                <div className="grid grid-cols-2 gap-4">
                  {[
                    { value: 'siswa', label: 'Student', emoji: '👨‍🎓', desc: 'Study and tasks' },
                    { value: 'guru', label: 'Teacher', emoji: '👨‍🏫', desc: 'Grade and manage' },
                  ].map(role => (
                    <label key={role.value}
                      className={`relative flex flex-col items-center text-center p-4 rounded-2xl border-2 cursor-pointer transition-all ${
                        formData.role === role.value
                          ? 'border-emerald-500 bg-emerald-50/50 shadow-md transform -translate-y-1'
                          : 'border-slate-100 bg-white hover:border-slate-200 hover:bg-slate-50'
                      }`}>
                      <input type="radio" name="role" value={role.value}
                        checked={formData.role === role.value}
                        onChange={handleChange}
                        className="sr-only" />
                      <span className="text-2xl mb-2">{role.emoji}</span>
                      <span className={`text-sm font-black uppercase tracking-tight ${formData.role === role.value ? 'text-emerald-700' : 'text-slate-700'}`}>
                        {role.label}
                      </span>
                      <span className="text-[10px] font-bold text-slate-400 mt-1 uppercase tracking-tighter opacity-70">{role.desc}</span>
                      {formData.role === role.value && (
                        <div className="absolute top-2 right-2 w-5 h-5 bg-emerald-500 rounded-full flex items-center justify-center text-white shadow-lg shadow-emerald-500/30">
                          <Check size={12} strokeWidth={4} />
                        </div>
                      )}
                    </label>
                  ))}
                </div>
              </div>

              <div className="pt-4">
                <Button
                  type="submit"
                  isLoading={isLoading}
                  fullWidth
                  size="lg"
                   className="bg-emerald-600! hover:bg-emerald-700! font-black tracking-wide shadow-emerald-500/20"
                >
                  Create Account →
                </Button>
              </div>
            </form>
          </div>

          <p className="text-center mt-10 text-sm text-slate-500 font-medium">
            Already have an account?{' '}
            <Link href="/login" className="font-black text-emerald-600 hover:text-emerald-500 hover:underline transition-all">
              Sign in instead
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
