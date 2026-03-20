"use client";

import React, { useEffect, useState, useRef } from 'react';
import api from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Plus, FileText, Trash2, Edit, X, Upload, Clock, BookOpen, Star, ChevronDown, Check } from 'lucide-react';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';

// ── Types ────────────────────────────────────────────────────────────────────
// ... (omitting types for brevity in replace_file_content if possible, but I'll keep them to avoid errors)
interface TaskItem {
  Id?: string;
  id?: string;
  name_task?: string;
  Name_Task?: string;
  file_task?: string;
  File_Task?: string;
  date_task?: string;
  Date_Task?: string;
  student_id?: string;
  Student_Id?: string;
  mapel_task?: string;
  Mapel_Task?: string;
  MapelTask?: string;
  created_at?: string;
  Created_at?: string;
  updated_at?: string;
  Updated_at?: string;
}

interface TaskWithStudent {
  Id: string;
  Name_Task: string;
  File_Task: string;
  Date_Task: string;
  Student_Id: string;
  MapelTask: string;
  Students?: { Full_Name?: string; Kelas?: string; Jurusan?: string };
}

interface GradeItem {
  Id: string;
  Task_Id: string;
  Grades: number;
  Keterangan: string;
  Tanggal: string;
}

// ── Helpers ──────────────────────────────────────────────────────────────────

const formatDate = (d: string) => {
  if (!d) return '—';
  try { return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }); }
  catch { return d; }
};

const getFileUrl = (filePath: string) => {
  if (!filePath) return null;
  const normalized = filePath.replace(/\\/g, '/');
  const filename = normalized.split('/').pop() || filePath;
  return `http://localhost:8080/api/v1/taskfile/${filename}`;
};

// ── Main Component ────────────────────────────────────────────────────────────

export default function TasksPage() {
  const { user } = useAuth();
  const isGuru = user?.role === 'guru';
  const isSiswa = user?.role === 'siswa';

  // Siswa state
  const [tasks, setTasks] = useState<TaskItem[]>([]);
  const [studentId, setStudentId] = useState<string | null>(null);

  // Guru state
  const [guruTasks, setGuruTasks] = useState<TaskWithStudent[]>([]);
  const [grades, setGrades] = useState<GradeItem[]>([]);
  const [gradeModal, setGradeModal] = useState<{ open: boolean; taskId: string; taskName: string }>({
    open: false, taskId: '', taskName: '',
  });
  const [gradeForm, setGradeForm] = useState({ grades: '', keterangan: '', tanggal: '' });
  const [gradeLoading, setGradeLoading] = useState(false);

  // Shared state
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({ name_task: '', mapel_task: '' });
  const [fileSelected, setFileSelected] = useState<File | null>(null);
  const [formLoading, setFormLoading] = useState(false);
  const [error, setError] = useState('');
  const fileRef = useRef<HTMLInputElement>(null);

  // ── Data Fetching ───────────────────────────────────────────────────────────

  const fetchSiswaData = async () => {
    if (!isSiswa) return;
    setLoading(true);
    try {
      const studentRes = await api.get('/students');
      const studentData = studentRes.data?.data;
      const sid = studentData?.id || studentData?.Id;
      if (!sid) {
        setLoading(false);
        return;
      }
      setStudentId(sid);
      const tasksRes = await api.get(`/tasks/student/${sid}`);
      setTasks(tasksRes.data?.data || []);
    } catch (err: any) {
      if (err?.code === 'ERR_CANCELED' || err?.name === 'CanceledError') return;
      console.error('Failed to fetch student tasks', err);
      setTasks([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchGuruData = async () => {
    if (!isGuru) return;
    setLoading(true);
    try {
      const [tasksRes, gradesRes] = await Promise.all([
        api.get('/tasks'),
        api.get('/grades'),
      ]);
      setGuruTasks(tasksRes.data?.data || []);
      setGrades(gradesRes.data?.data || []);
    } catch (err: any) {
      if (err?.code === 'ERR_CANCELED' || err?.name === 'CanceledError') return;
      console.error('Failed to fetch guru tasks', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isSiswa) fetchSiswaData();
    if (isGuru) fetchGuruData();
  }, [isSiswa, isGuru]);

  // ── Siswa CRUD ──────────────────────────────────────────────────────────────

  const openCreateModal = () => {
    setEditingId(null);
    setFormData({ name_task: '', mapel_task: '' });
    setFileSelected(null);
    setError('');
    setIsModalOpen(true);
  };

  const openEditModal = (task: TaskItem) => {
    const id = task.Id || task.id || '';
    setEditingId(id);
    setFormData({
      name_task: task.Name_Task || task.name_task || '',
      mapel_task: task.MapelTask || task.Mapel_Task || task.mapel_task || ''
    });
    setFileSelected(null);
    setError('');
    setIsModalOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingId && !fileSelected) {
      setError('Please upload a file for the task.');
      return;
    }
    setFormLoading(true);
    setError('');
    try {
      if (editingId) {
        const fd = new FormData();
        if (formData.name_task) fd.append('name_task', formData.name_task);
        if (formData.mapel_task) fd.append('mapel_task', formData.mapel_task);
        if (fileSelected) fd.append('file_task', fileSelected);
        await api.patch(`/task/${editingId}`, fd, { headers: { 'Content-Type': 'multipart/form-data' } });
      } else {
        const fd = new FormData();
        fd.append('name_task', formData.name_task);
        fd.append('mapel_task', formData.mapel_task);
        if (studentId) fd.append('student_id', studentId);
        if (fileSelected) fd.append('file_task', fileSelected);
        await api.post('/task', fd, { headers: { 'Content-Type': 'multipart/form-data' } });
      }
      setIsModalOpen(false);
      fetchSiswaData();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to save task. Check file type (jpg/png/pdf allowed).');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this task?')) return;
    try {
      await api.delete(`/task/${id}`);
      setTasks(prev => prev.filter(t => (t.Id || t.id) !== id));
    } catch (err: any) {
      alert(err.response?.data?.message || 'Failed to delete task.');
    }
  };

  // ── Guru Grading ────────────────────────────────────────────────────────────

  const openGradeModal = (taskId: string, taskName: string) => {
    setGradeForm({ grades: '', keterangan: '', tanggal: new Date().toISOString().split('T')[0] });
    setGradeModal({ open: true, taskId, taskName });
    setError('');
  };

  const handleGradeSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setGradeLoading(true);
    setError('');
    try {
      await api.post('/grades', {
        task_id: gradeModal.taskId,
        grades: parseInt(gradeForm.grades),
        keterangan: gradeForm.keterangan,
        tanggal: gradeForm.tanggal,
      });
      setGradeModal({ open: false, taskId: '', taskName: '' });
      fetchGuruData();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to submit grade.');
    } finally {
      setGradeLoading(false);
    }
  };

  const getGradeForTask = (taskId: string) => grades.find(g => g.Task_Id === taskId);

  // ── Render: Siswa View ──────────────────────────────────────────────────────

  if (isSiswa) {
    return (
      <div className="space-y-8 animate-fade-in-up">
        {/* Header */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
          <div className="space-y-2">
            <h1 className="text-2xl font-bold text-white leading-tight flex items-center gap-3 uppercase tracking-tight">
              <BookOpen size={24} className="text-indigo-400" /> Tugas Saya
            </h1>
            <p className="text-slate-400 text-sm leading-relaxed font-medium">Upload dan kelola tugas-tugas kamu di sini.</p>
          </div>
          <Button
            onClick={openCreateModal}
            icon={<Plus size={18} />}
            size="md"
            className="font-black tracking-wide"
          >
            Upload Tugas Baru
          </Button>
        </div>

        {/* No student profile warning */}
        {!loading && !studentId && (
          <div className="glass border border-amber-500/20 bg-amber-500/5 text-amber-200 rounded-2xl p-6 text-sm font-bold leading-relaxed flex items-center gap-4 animate-pulse">
            <span className="text-xl">⚠️</span>
            <span>Kamu belum terdaftar sebagai siswa. Daftar dulu di menu <strong className="text-white underline underline-offset-4">Students</strong>.</span>
          </div>
        )}

        {/* Task list */}
        {loading ? (
          <div className="glass rounded-3xl p-16 text-center space-y-4">
            <div className="animate-spin w-10 h-10 border-4 border-indigo-500 border-t-transparent rounded-full mx-auto" />
            <p className="text-slate-400 text-sm font-bold uppercase tracking-widest">Loading tugas...</p>
          </div>
        ) : tasks.length === 0 && studentId ? (
          <div className="glass border-2 border-dashed border-white/5 rounded-3xl p-20 text-center space-y-6">
            <div className="w-24 h-24 mx-auto rounded-3xl bg-white/5 flex items-center justify-center text-slate-600 shadow-inner">
              <FileText size={48} />
            </div>
            <div className="space-y-2">
              <h3 className="text-xl font-black text-white leading-tight uppercase">Belum Ada Tugas</h3>
              <p className="text-slate-500 text-sm font-medium max-w-xs mx-auto">Klik tombol &quot;Upload Tugas Baru&quot; untuk mengumpulkan tugas pertama kamu.</p>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6 md:gap-8">
            {tasks.map((task) => {
              const taskId = task.Id || task.id || '';
              const nameTask = task.Name_Task || task.name_task || 'Untitled';
              const mapelTask = task.MapelTask || task.Mapel_Task || task.mapel_task || 'Umum';
              const dateTask = task.Date_Task || task.date_task || '';
              const fileTask = task.File_Task || task.file_task || '';
              return (
                <div key={taskId} className="glass rounded-2xl overflow-hidden hover:shadow-2xl transition-all group flex flex-col border border-white/5 hover:border-white/10">
                  <div className="h-1.5 w-full bg-linear-to-r from-indigo-600 to-purple-600" />
                  <div className="p-8 flex flex-col flex-1 space-y-5">
                    <div className="flex items-start justify-between gap-4">
                      <span className="badge badge-indigo py-1.5 px-4 font-black uppercase tracking-wider text-[10px]">
                        {mapelTask}
                      </span>
                      <div className="flex items-center gap-3 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button
                          onClick={() => openEditModal(task)}
                          className="p-2 text-slate-400 hover:text-white hover:bg-white/10 rounded-xl transition-all" title="Edit"
                        >
                          <Edit size={16} />
                        </button>
                        <button
                          onClick={() => handleDelete(taskId)}
                          className="p-2 text-slate-400 hover:text-red-400 hover:bg-red-400/10 rounded-xl transition-all" title="Delete"
                        >
                          <Trash2 size={16} />
                        </button>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <h3 className="text-lg font-black text-white leading-tight line-clamp-2 min-h-14 uppercase tracking-tight">{nameTask}</h3>
                      <p className="text-[10px] text-slate-500 font-black uppercase tracking-widest flex items-center gap-2 leading-none">
                        <Clock size={14} className="text-indigo-500" /> {formatDate(dateTask)}
                      </p>
                    </div>

                    {fileTask && (
                      <div className="pt-2">
                        <a
                          href={getFileUrl(fileTask) || '#'}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="w-full inline-flex items-center justify-center gap-3 text-[10px] font-black uppercase tracking-[0.2em] text-indigo-400 bg-indigo-400/5 hover:bg-indigo-400/10 border border-indigo-400/20 py-3 rounded-xl transition-all"
                        >
                          <FileText size={14} /> Lihat File Tugas
                        </a>
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* Create/Edit Modal */}
        {isModalOpen && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-6">
            <div className="absolute inset-0 bg-black/80 backdrop-blur-md" onClick={() => setIsModalOpen(false)} />
            <div className="relative glass rounded-3xl shadow-2xl w-full max-w-lg overflow-hidden border border-white/10 animate-fade-in-up" onClick={e => e.stopPropagation()}>
              <div className="h-2 w-full bg-linear-to-r from-indigo-500 to-purple-500" />
              <div className="flex items-center justify-between px-8 py-6 border-b border-white/5 bg-white/2">
                <div className="space-y-1">
                  <h2 className="text-xl font-black text-white leading-tight uppercase tracking-tight">{editingId ? 'Edit Tugas' : 'Upload Tugas Baru'}</h2>
                  <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest leading-none">Isi informasi tugas kamu dengan lengkap.</p>
                </div>
                <button onClick={() => setIsModalOpen(false)} className="text-slate-400 hover:text-white hover:bg-white/10 p-2.5 rounded-xl transition-colors">
                  <X size={20} />
                </button>
              </div>

              <form onSubmit={handleSubmit} className="p-8 space-y-6">
                {error && (
                  <div className="p-4 bg-red-500/10 border border-red-500/20 text-red-500 rounded-2xl text-sm font-semibold leading-normal flex items-center gap-3 animate-head-shake">
                    <span className="text-lg">⚠</span> {error}
                  </div>
                )}

                <Input
                  label="Nama Tugas"
                  type="text"
                  value={formData.name_task}
                  onChange={e => setFormData({ ...formData, name_task: e.target.value })}
                  placeholder="Contoh: Tugas Matematika Bab 3"
                  required
                />

                <Input
                  label="Mata Pelajaran"
                  type="text"
                  value={formData.mapel_task}
                  onChange={e => setFormData({ ...formData, mapel_task: e.target.value })}
                  placeholder="Contoh: Matematika"
                  required
                />

                <div className="space-y-2">
                  <label className="form-label">
                    File Tugas {editingId ? '(opsional)' : '*'}
                  </label>
                  <input
                    ref={fileRef}
                    type="file"
                    accept=".jpg,.jpeg,.png,.pdf"
                    className="sr-only"
                    onChange={e => setFileSelected(e.target.files?.[0] || null)}
                  />
                  <button
                    type="button"
                    onClick={() => fileRef.current?.click()}
                    className={`w-full flex flex-col items-center justify-center gap-4 p-8 border-2 border-dashed rounded-3xl transition-all group ${
                      fileSelected ? 'border-emerald-500/40 bg-emerald-500/5' : 'border-white/10 hover:border-indigo-500/50 hover:bg-indigo-500/5'
                    }`}
                  >
                    <div className={`w-14 h-14 rounded-2xl flex items-center justify-center transition-all ${
                      fileSelected ? 'bg-emerald-500/20 text-emerald-400' : 'bg-white/5 text-slate-400 group-hover:text-indigo-400'
                    }`}>
                      {fileSelected ? <Check size={28} /> : <Upload size={28} />}
                    </div>
                    <div className="text-center">
                      <p className={`text-sm font-black uppercase tracking-tight transition-colors ${fileSelected ? 'text-emerald-400' : 'text-slate-300 group-hover:text-white'}`}>
                        {fileSelected ? fileSelected.name : 'Klik untuk pilih file'}
                      </p>
                      <p className="text-[10px] text-slate-500 mt-2 uppercase tracking-[0.2em] font-black opacity-60">JPG, PNG, PDF • MAX 10MB</p>
                    </div>
                  </button>
                </div>

                <div className="pt-6 flex justify-end gap-4 border-t border-white/5">
                  <Button
                    type="button"
                    variant="secondary"
                    onClick={() => setIsModalOpen(false)}
                    className="font-bold tracking-tight"
                  >
                    Batal
                  </Button>
                  <Button
                    type="submit"
                    isLoading={formLoading}
                    className="font-black tracking-wide"
                  >
                    {editingId ? 'Simpan Perubahan' : 'Upload Tugas'}
                  </Button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    );
  }

  // ── Render: Guru View ───────────────────────────────────────────────────────

  return (
    <div className="space-y-8 animate-fade-in-up">
      {/* Header */}
      <div className="space-y-2">
        <h1 className="text-2xl font-black text-white leading-tight flex items-center gap-3 uppercase tracking-tight">
          <BookOpen size={24} className="text-emerald-400" /> Tugas Siswa
        </h1>
        <p className="text-slate-400 text-sm font-medium leading-relaxed">Lihat tugas yang dikumpulkan siswa dan berikan penilaian langsung.</p>
      </div>

      {loading ? (
        <div className="glass rounded-3xl p-16 text-center space-y-4">
          <div className="animate-spin w-10 h-10 border-4 border-emerald-500 border-t-transparent rounded-full mx-auto" />
          <p className="text-slate-400 text-sm font-bold uppercase tracking-widest">Loading data tugas...</p>
        </div>
      ) : guruTasks.length === 0 ? (
        <div className="glass border-2 border-dashed border-white/5 rounded-3xl p-20 text-center space-y-6">
          <div className="w-24 h-24 mx-auto rounded-3xl bg-white/5 flex items-center justify-center text-slate-600 shadow-inner">
            <FileText size={48} />
          </div>
          <div className="space-y-2">
            <h3 className="text-xl font-black text-white leading-tight uppercase">Belum Ada Tugas</h3>
            <p className="text-slate-500 text-sm font-medium max-w-xs mx-auto">Siswa Anda belum mulai mengumpulkan tugas untuk saat ini.</p>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 md:gap-8">
          {guruTasks.map((task) => {
            const id = task.Id;
            const existingGrade = getGradeForTask(id);
            return (
              <div key={id} className="glass rounded-3xl overflow-hidden hover:shadow-2xl transition-all flex flex-col border border-white/5 hover:border-white/10 group">
                <div className="h-2 w-full" style={{ background: existingGrade ? 'linear-gradient(90deg, #10b981, #059669)' : 'linear-gradient(90deg, #f59e0b, #eab308)' }} />
                <div className="p-8 flex flex-col flex-1 space-y-5">
                  <div className="flex items-start justify-between gap-4">
                    <span className="badge badge-indigo py-1.5 px-4 font-black uppercase tracking-wider text-[10px]">
                      {task.MapelTask || 'Umum'}
                    </span>
                    {existingGrade && (
                      <div className="flex items-center gap-2 bg-emerald-500/10 text-emerald-400 px-3 py-1.5 rounded-full border border-emerald-500/20 shadow-inner">
                        <Star size={12} fill="currentColor" />
                        <span className="text-sm font-black leading-none">{existingGrade.Grades}</span>
                      </div>
                    )}
                  </div>

                  <div className="space-y-3">
                    <h3 className="text-lg font-black text-white leading-tight line-clamp-2 min-h-14 group-hover:text-emerald-400 transition-colors uppercase tracking-tight">{task.Name_Task}</h3>
                    <div className="space-y-2 pt-1">
                      <p className="text-[11px] text-slate-200 font-black flex items-center gap-2 leading-none uppercase tracking-widest">
                        <span className="opacity-70 text-base">👤</span> {task.Students?.Full_Name || 'Unknown'}
                      </p>
                      <p className="text-[10px] text-slate-600 font-bold flex items-center gap-2 leading-none uppercase tracking-tighter">
                        <Clock size={14} className="text-slate-700" /> {formatDate(task.Date_Task)}
                      </p>
                    </div>
                  </div>

                  {existingGrade && (
                    <div className="bg-white/5 p-4 rounded-2xl border border-white/5">
                      <p className="text-[11px] text-slate-400 italic font-medium leading-relaxed line-clamp-2">
                        &ldquo;{existingGrade.Keterangan}&rdquo;
                      </p>
                    </div>
                  )}

                  <div className="mt-auto pt-4 flex flex-col sm:flex-row items-center gap-3">
                    {task.File_Task && (
                      <a
                        href={getFileUrl(task.File_Task) || '#'}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="w-full sm:flex-1 inline-flex items-center justify-center gap-2 text-[10px] font-black uppercase tracking-widest text-indigo-400 hover:text-white bg-indigo-400/5 hover:bg-indigo-600 border border-indigo-400/20 py-3 rounded-xl transition-all"
                      >
                        <FileText size={14} /> File
                      </a>
                    )}
                    <button
                      onClick={() => openGradeModal(id, task.Name_Task)}
                      className="w-full sm:flex-1 flex items-center justify-center gap-2 text-[10px] font-black uppercase tracking-widest text-white py-3 rounded-xl transition-all active:scale-95 shadow-lg shadow-black/20"
                      style={{ background: existingGrade ? 'linear-gradient(135deg, #10b981, #059669)' : 'linear-gradient(135deg, #4f46e5, #7c3aed)' }}
                    >
                      <Star size={14} /> {existingGrade ? 'Update' : 'Nilai'}
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Grade Modal */}
      {gradeModal.open && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-6">
          <div className="absolute inset-0 bg-black/80 backdrop-blur-md" onClick={() => setGradeModal({ open: false, taskId: '', taskName: '' })} />
          <div className="relative glass rounded-3xl shadow-2xl w-full max-w-md overflow-hidden border border-white/10 animate-fade-in-up" onClick={e => e.stopPropagation()}>
            <div className="h-2 w-full bg-linear-to-r from-emerald-500 to-indigo-500" />
            <div className="flex items-center justify-between px-8 py-6 border-b border-white/5 bg-white/2">
              <div className="space-y-1">
                <h2 className="text-xl font-black text-white leading-tight uppercase tracking-tight">Penilaian Tugas</h2>
                <p className="text-[10px] text-emerald-400 font-black uppercase tracking-widest leading-none truncate max-w-[250px]">{gradeModal.taskName}</p>
              </div>
              <button onClick={() => setGradeModal({ open: false, taskId: '', taskName: '' })} className="text-slate-400 hover:text-white hover:bg-white/10 p-2.5 rounded-xl transition-colors">
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleGradeSubmit} className="p-8 space-y-6">
              {error && (
                <div className="p-4 bg-red-500/10 border border-red-500/20 text-red-500 rounded-2xl text-sm font-semibold leading-normal flex items-center gap-3 animate-head-shake">
                  <span className="text-lg">⚠</span> {error}
                </div>
              )}

              <Input
                label="Nilai (0-100)"
                type="number"
                min="0"
                max="100"
                value={gradeForm.grades}
                onChange={e => setGradeForm({ ...gradeForm, grades: e.target.value })}
                required
                placeholder="Contoh: 85"
                className="text-lg font-black"
                containerClassName="animate-fade-in-up md:animation-delay-100"
              />

              <Input
                label="Keterangan / Feedback"
                as="textarea"
                value={gradeForm.keterangan}
                onChange={e => setGradeForm({ ...gradeForm, keterangan: e.target.value })}
                required
                rows={4}
                placeholder="Berikan feedback untuk siswa..."
                containerClassName="animate-fade-in-up md:animation-delay-200"
              />

              <Input
                label="Tanggal Penilaian"
                type="date"
                value={gradeForm.tanggal}
                onChange={e => setGradeForm({ ...gradeForm, tanggal: e.target.value })}
                required
                containerClassName="animate-fade-in-up md:animation-delay-300"
              />

              <div className="pt-6 flex justify-end gap-4 border-t border-white/5 animate-fade-in-up md:animation-delay-400">
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => setGradeModal({ open: false, taskId: '', taskName: '' })}
                  className="font-bold tracking-tight"
                >
                  Batal
                </Button>
                  <Button
                    type="submit"
                    isLoading={gradeLoading}
                    variant="success"
                    className="bg-emerald-600! hover:bg-emerald-700! font-black tracking-wide shadow-emerald-500/20"
                  >
                  Simpan Nilai
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
