"use client";

import React, { useEffect, useState } from 'react';
import api from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Award, TrendingUp, Plus, Edit, Trash2, BookOpen, Star, Calendar, Clock, ChevronRight } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';

// ── Main Component ────────────────────────────────────────────────────────────

export default function GradesPage() {
  const { user } = useAuth();
  const isGuru = user?.role === 'guru';
  const isSiswa = user?.role === 'siswa';
  const [grades, setGrades] = useState<any[]>([]);
  const [tasks, setTasks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  // Form State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({ 
    task_id: '', 
    tanggal: '', 
    keterangan: '', 
    grades: '' 
  });
  const [formLoading, setFormLoading] = useState(false);
  const [error, setError] = useState('');

  const fetchGradesData = async () => {
    try {
      setLoading(true);
      if (isGuru) {
        const [res, tasksRes] = await Promise.all([
          api.get('/grades'),
          api.get('/tasks')
        ]);
        setGrades(res.data?.data || []);
        setTasks(tasksRes.data?.data || []);
      } else if (isSiswa) {
        const studentRes = await api.get('/students');
        const sid = studentRes.data?.data?.id || studentRes.data?.data?.Id;
        if (sid) {
          const gradesRes = await api.get(`/grades/student/${sid}`);
          setGrades(gradesRes.data?.data || []);
          const tasksRes = await api.get(`/tasks/student/${sid}`);
          setTasks(tasksRes.data?.data || []);
        } else {
          setGrades([]);
        }
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGradesData();
  }, [isGuru, isSiswa]);

  const openAddModal = () => {
    setEditingId(null);
    setFormData({ 
      task_id: tasks.length > 0 ? (tasks[0].Id || tasks[0].id) : '', 
      tanggal: new Date().toISOString().split('T')[0], 
      keterangan: 'Passed with excellence', 
      grades: '100' 
    });
    setError('');
    setIsModalOpen(true);
  };

  const openEditModal = (grade: any) => {
    setEditingId(grade.Id || grade.id);
    setFormData({ 
      task_id: grade.Task_Id || grade.task_id || '', 
      tanggal: grade.Tanggal ? new Date(grade.Tanggal).toISOString().split('T')[0] : '', 
      keterangan: grade.Keterangan || grade.keterangan || '', 
      grades: grade.Grades || grade.grades || '0' 
    });
    setError('');
    setIsModalOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormLoading(true);
    setError('');
    try {
      const payload = {
        ...formData,
        tanggal: new Date(formData.tanggal).toISOString(),
        grades: parseInt(formData.grades.toString())
      };

      if (editingId) {
        await api.patch(`/grades/${editingId}`, payload);
      } else {
        await api.post('/grades', payload);
      }
      setIsModalOpen(false);
      fetchGradesData();
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.message || 'Failed to save grade.');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this grade record?')) return;
    try {
      await api.delete(`/grades/${id}`);
      setGrades(grades.filter(g => (g.Id || g.id) !== id));
    } catch (err) {
      console.error(err);
      alert('Failed to delete grade.');
    }
  };

  const getPerformanceTier = (score: number) => {
    if (score >= 90) return { label: 'A+', color: 'text-emerald-400', bg: 'bg-emerald-500/10' };
    if (score >= 80) return { label: 'A', color: 'text-indigo-400', bg: 'bg-indigo-500/10' };
    if (score >= 70) return { label: 'B', color: 'text-blue-400', bg: 'bg-blue-500/10' };
    return { label: 'C', color: 'text-amber-400', bg: 'bg-amber-500/10' };
  };

  const averageGrade = grades.length > 0 
    ? Math.round(grades.reduce((acc, curr) => acc + (curr.Grades || curr.grades || 0), 0) / grades.length) 
    : 0;

  return (
    <div className="space-y-8 animate-fade-in-up">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
        <div className="space-y-2">
          <h1 className="text-2xl font-black text-white leading-tight flex items-center gap-3 uppercase tracking-tight">
            <Award size={24} className="text-indigo-400" /> Hasil Penilaian
          </h1>
          <p className="text-slate-400 text-sm font-medium leading-relaxed">Pantau evaluasi belajar dan ringkasan nilai akademik kamu.</p>
        </div>
        {isGuru && (
          <Button 
            onClick={openAddModal} 
            icon={<Plus size={18} />}
            className="font-black tracking-wide"
          >
            Submit Grades
          </Button>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 md:gap-10">
        <div className="col-span-1 lg:col-span-2 space-y-6">
          {loading ? (
             <div className="glass rounded-3xl p-16 text-center space-y-4">
                <div className="animate-spin w-10 h-10 border-4 border-indigo-500 border-t-transparent rounded-full mx-auto" />
                <p className="text-slate-400 text-sm font-bold uppercase tracking-widest">Loading grades data...</p>
             </div>
          ) : grades.length === 0 ? (
            <div className="glass rounded-3xl p-20 text-center space-y-6 border-2 border-dashed border-white/5">
              <div className="w-24 h-24 mx-auto rounded-3xl bg-white/5 flex items-center justify-center text-slate-600 shadow-inner">
                <Award size={48} />
              </div>
              <div className="space-y-2">
                <h3 className="text-xl font-black text-white leading-tight uppercase">Belum Ada Nilai</h3>
                <p className="text-slate-500 text-sm font-medium max-w-xs mx-auto">Evaluasi belajar Anda belum tersedia. Pastikan tugas sudah dikumpulkan dan dinilai.</p>
              </div>
            </div>
          ) : (
            grades.map((grade, i) => {
              const id = grade.Id || grade.id || i;
              const title = grade.Task_Name || grade.task_name || 'Assignment Task';
              const mapel = grade.Mapel_Task || grade.mapel_task || 'General subject';
              const keterangan = grade.Keterangan || grade.keterangan || 'No additional feedback provided.';
              const score = grade.Grades || grade.grades || 0;
              const date = grade.Tanggal || grade.tanggal || grade.Created_at;
              const tier = getPerformanceTier(score);

              return (
              <div key={id} className="glass group rounded-3xl overflow-hidden hover:shadow-2xl transition-all border border-white/5 hover:border-white/10 flex flex-col sm:flex-row items-center gap-6 p-6 md:p-8">
                <div className={`w-20 h-20 rounded-2xl flex flex-col items-center justify-center shrink-0 border-2 border-white/5 shadow-inner group-hover:scale-105 transition-transform ${tier.bg}`}>
                  <span className={`text-2xl font-black leading-none ${tier.color}`}>{score}</span>
                  <span className="text-[9px] font-black uppercase tracking-tighter opacity-50 mt-1">PTS</span>
                </div>
                <div className="flex-1 text-center sm:text-left min-w-0 space-y-3">
                  <div className="space-y-1">
                    <h4 className="font-black text-white text-lg truncate leading-tight group-hover:text-indigo-400 transition-colors uppercase tracking-tight">{title}</h4>
                    <p className="text-[10px] font-black uppercase tracking-widest text-indigo-500 leading-none flex items-center justify-center sm:justify-start gap-2">
                      <BookOpen size={12} /> {mapel}
                    </p>
                  </div>
                  <div className="bg-white/5 px-4 py-3 rounded-2xl border border-white/5">
                    <p className="text-xs text-slate-400 leading-relaxed italic line-clamp-2">&ldquo;{keterangan}&rdquo;</p>
                  </div>
                </div>
                <div className="flex flex-col items-center sm:items-end justify-between self-stretch py-1">
                  <div className="text-[10px] font-black uppercase tracking-widest text-slate-500 bg-white/5 px-3 py-1.5 rounded-full leading-none flex items-center gap-2">
                    <Calendar size={12} className="text-slate-700" />
                    {new Date(date || Date.now()).toLocaleDateString('id-ID', {day:'numeric', month:'short'})}
                  </div>
                  {isGuru && (
                    <div className="flex justify-end gap-3 opacity-0 group-hover:opacity-100 transition-all translate-y-2 group-hover:translate-y-0 pt-4 sm:pt-0">
                      <button 
                        className="p-2.5 text-slate-500 hover:text-white hover:bg-white/10 rounded-xl transition-all"
                        onClick={() => openEditModal(grade)}
                        title="Edit"
                      >
                        <Edit size={16} />
                      </button>
                      <button 
                        className="p-2.5 text-slate-500 hover:text-red-400 hover:bg-red-400/10 rounded-xl transition-all"
                        onClick={() => handleDelete(id)}
                        title="Delete"
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  )}
                </div>
              </div>
            )})
          )}
        </div>

        <div className="col-span-1">
          <div className="glass rounded-3xl p-8 space-y-8 sticky top-8 border border-white/5 shadow-2xl overflow-hidden bg-linear-to-b from-indigo-950/50 to-slate-950/80">
            <div className="absolute -top-10 -right-10 w-40 h-40 bg-indigo-600 rounded-full blur-[80px] opacity-20 pointer-events-none" />
            
            <div className="space-y-4 relative">
              <h3 className="text-lg font-black text-white flex items-center gap-3 leading-tight uppercase tracking-wider">
                <TrendingUp size={20} className="text-indigo-400" /> Ringkasan Performa
              </h3>
              <p className="text-xs text-slate-500 font-medium leading-relaxed">Analisis performa belajar Anda berdasarkan evaluasi tugas yang telah dinilai.</p>
            </div>
            
            <div className="space-y-4 pt-4 border-t border-white/10">
              <div className="bg-white/2 border border-white/5 p-5 rounded-2xl flex justify-between items-center group hover:bg-white/5 transition-colors">
                <span className="text-[11px] font-black uppercase tracking-widest text-slate-500">Total Evaluated</span>
                <span className="text-indigo-400 font-black text-sm">{grades.length} Tasks</span>
              </div>
              <div className="bg-white/2 border border-white/5 p-5 rounded-2xl flex justify-between items-center group hover:bg-white/5 transition-colors">
                <span className="text-[11px] font-black uppercase tracking-widest text-slate-500">Tracking Status</span>
                <span className="text-emerald-400 font-black text-sm flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" /> ACTIVE
                </span>
              </div>
            </div>

            <div className="pt-4">
              <div className="p-6 rounded-3xl bg-linear-to-br from-indigo-600/20 to-purple-600/20 border border-indigo-500/20 text-center shadow-inner">
                 <p className="text-[10px] font-black text-indigo-400 uppercase tracking-[0.25em] mb-2">GPA ESTIMATION</p>
                 <p className="text-4xl font-black text-white leading-none tracking-tighter">
                   {averageGrade >= 90 ? 'A+' : averageGrade >= 80 ? 'A' : averageGrade >= 70 ? 'B' : 'C'}
                 </p>
                 <div className="mt-4 h-1.5 w-full bg-white/5 rounded-full overflow-hidden">
                    <div className="h-full bg-indigo-500 rounded-full shadow-lg shadow-indigo-500/50" style={{ width: `${averageGrade}%` }} />
                 </div>
                 <p className="text-[10px] font-black text-slate-500 mt-2 uppercase tracking-widest">Score Average: {averageGrade}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        title={editingId ? "Edit Grade" : "Submit Grade"}
      >
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="p-4 bg-red-500/10 border border-red-500/20 text-red-500 rounded-2xl text-sm font-semibold leading-normal flex items-center gap-3 animate-head-shake">
              <span className="text-lg">⚠</span> {error}
            </div>
          )}

          <Input
            label="Tugas / Assignment"
            as="select"
            value={formData.task_id}
            onChange={e => setFormData({...formData, task_id: e.target.value})}
            required
            icon={<BookOpen size={18} />}
          >
            <option value="" disabled>Select Task</option>
            {tasks.map(t => (
              <option key={t.Id || t.id} value={t.Id || t.id}>{t.Name_Task || t.name_task || 'Unknown Task'}</option>
            ))}
          </Input>
          
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
            <Input 
              label="Tanggal Penilaian" 
              type="date"
              value={formData.tanggal} 
              onChange={e => setFormData({...formData, tanggal: e.target.value})} 
              required 
            />
            <Input 
              label="Nilai (0-100)" 
              type="number"
              min="0"
              max="100"
              value={formData.grades} 
              onChange={e => setFormData({...formData, grades: e.target.value})} 
              required 
              placeholder="Contoh: 95"
            />
          </div>

          <Input
            label="Feedback Akademik"
            as="textarea"
            value={formData.keterangan}
            onChange={e => setFormData({...formData, keterangan: e.target.value})}
            required
            rows={4}
            placeholder="Berikan feedback akademik untuk siswa..."
          />
          
          <div className="pt-6 flex justify-end gap-4 border-t border-white/5">
            <Button 
              type="button" 
              variant="secondary"
              onClick={() => setIsModalOpen(false)}
              className="font-bold tracking-tight"
            >
              Cancel
            </Button>
            <Button 
              type="submit" 
              isLoading={formLoading}
              className="font-black tracking-wide"
            >
              {editingId ? 'Simpan Perubahan' : 'Save Grade'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
