"use client";

import React, { useEffect, useState } from 'react';
import api from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Award, TrendingUp, Plus, Edit, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';

export default function GradesPage() {
  const { user } = useAuth();
  const isGuru = user?.role === 'guru';
  const isSiswa = user?.role === 'siswa';
  const [grades, setGrades] = useState<any[]>([]);
  const [tasks, setTasks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  // Form State based on PayloadsStudentGrade
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({ 
    task_id: '', 
    tanggal: '', 
    keterangan: '', 
    grades: '' 
  });
  const [formLoading, setFormLoading] = useState(false);

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
        // Fetch student data first to get sid
        const studentRes = await api.get('/students');
        const sid = studentRes.data?.data?.id || studentRes.data?.data?.Id;
        
        if (sid) {
          // Now fetch grades using the new student-specific endpoint
          const gradesRes = await api.get(`/grades/student/${sid}`);
          setGrades(gradesRes.data?.data || []);
          
          // Also fetch tasks if needed for context (though not strictly required for the list if using the new endpoint)
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
      keterangan: 'Passed', 
      grades: '100' 
    });
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
    setIsModalOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormLoading(true);
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
      alert(err.response?.data?.message || 'Failed to save grade.');
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

  return (
    <div className="space-y-8 animate-fade-in-up">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
        <div className="space-y-2">
          <h1 className="text-2xl font-bold text-white leading-tight flex items-center gap-3">
            <Award size={24} className="text-indigo-400" /> Grades & Scores
          </h1>
          <p className="text-slate-400 text-sm leading-relaxed">Student evaluation and performance metrics summary.</p>
        </div>
        {isGuru && (
          <button onClick={openAddModal} className="btn-primary shrink-0 transition-transform active:scale-95 shadow-lg shadow-indigo-500/20">
            <Plus size={16} /> Submit Grades
          </button>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 md:gap-10">
        <div className="col-span-1 lg:col-span-2 space-y-6">
          {loading ? (
             <div className="glass rounded-3xl p-16 text-center space-y-4">
                <div className="animate-spin w-10 h-10 border-4 border-indigo-500 border-t-transparent rounded-full mx-auto" />
                <p className="text-slate-400 text-sm leading-relaxed">Loading grades data...</p>
             </div>
          ) : grades.length === 0 ? (
            <div className="glass rounded-3xl p-20 text-center space-y-6 border-2 border-dashed border-white/5">
              <div className="w-20 h-20 mx-auto rounded-3xl bg-white/5 flex items-center justify-center text-slate-500">
                <Award size={48} />
              </div>
              <div className="space-y-2">
                <h3 className="text-xl font-bold text-white leading-tight">No grades available</h3>
                <p className="text-slate-400 text-sm leading-relaxed max-w-xs mx-auto">Evaluasi belajar Anda belum tersedia. Pastikan tugas sudah dikumpulkan.</p>
              </div>
            </div>
          ) : (
            grades.map((grade, i) => {
              const id = grade.Id || grade.id || i;
              const title = grade.Task_Name || grade.task_name || 'Assignment Task';
              const mapel = grade.Mapel_Task || grade.mapel_task || 'General subject';
              const keterangan = grade.Keterangan || grade.keterangan || 'No additional feedback provided.';
              const score = grade.Grades || grade.grades || '0';
              const date = grade.Tanggal || grade.tanggal || grade.Created_at;

              return (
              <div key={id} className="glass group rounded-2xl overflow-hidden hover:shadow-2xl transition-all border border-white/5 hover:border-white/10 flex flex-col sm:flex-row items-center gap-6 p-6 md:p-8">
                <div className="w-20 h-20 rounded-2xl flex flex-col items-center justify-center shrink-0 border-2 border-indigo-500/20 shadow-inner group-hover:scale-105 transition-transform" style={{ background: 'linear-gradient(135deg, rgba(79,70,229,0.1), rgba(124,58,237,0.1))' }}>
                  <span className="text-2xl font-black text-indigo-400 leading-none">{score}</span>
                  <span className="text-[9px] font-black uppercase tracking-tighter text-indigo-500/50 mt-1">PTS</span>
                </div>
                <div className="flex-1 text-center sm:text-left min-w-0 space-y-2">
                  <div className="space-y-0.5">
                    <h4 className="font-black text-white text-lg truncate leading-tight group-hover:text-indigo-400 transition-colors uppercase tracking-tight">{title}</h4>
                    <p className="text-[11px] font-black uppercase tracking-widest text-indigo-500 leading-none">{mapel}</p>
                  </div>
                  <div className="bg-white/5 px-4 py-3 rounded-xl border border-white/5">
                    <p className="text-xs text-slate-400 leading-relaxed italic line-clamp-2">&ldquo;{keterangan}&rdquo;</p>
                  </div>
                </div>
                <div className="flex flex-col items-center sm:items-end justify-between self-stretch py-1">
                  <div className="text-[10px] font-black uppercase tracking-widest text-slate-500 bg-white/5 px-3 py-1 rounded-full leading-none">
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
          <div className="glass rounded-3xl p-8 space-y-8 relative overflow-hidden sticky top-8" style={{ background: 'linear-gradient(180deg, #1e1b4b 0%, #0d1117 100%)', border: '1px solid rgba(255,255,255,0.05)' }}>
            <div className="absolute -top-10 -right-10 w-40 h-40 bg-indigo-600 rounded-full blur-[80px] opacity-20 pointer-events-none"></div>
            
            <div className="space-y-4 relative">
              <h3 className="text-lg font-black text-white flex items-center gap-3 leading-tight uppercase tracking-wider">
                <TrendingUp size={20} className="text-indigo-400" /> Quick Stats
              </h3>
              <p className="text-xs text-slate-500 font-medium leading-relaxed">Analisis performa belajar Anda berdasarkan evaluasi tugas yang telah dinilai oleh guru mata pelajaran.</p>
            </div>
            
            <div className="space-y-4 pt-4 border-t border-white/5">
              <div className="bg-white/[.02] border border-white/5 p-5 rounded-2xl flex justify-between items-center group hover:bg-white/5 transition-colors">
                <span className="text-[11px] font-black uppercase tracking-widest text-slate-400">Total Evaluated</span>
                <span className="text-indigo-400 font-black text-sm">{grades.length} Tasks</span>
              </div>
              <div className="bg-white/[.02] border border-white/5 p-5 rounded-2xl flex justify-between items-center group hover:bg-white/5 transition-colors">
                <span className="text-[11px] font-black uppercase tracking-widest text-slate-400">Tracking Status</span>
                <span className="text-emerald-400 font-black text-sm flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" /> ACTIVE
                </span>
              </div>
            </div>

            <div className="pt-4">
              <div className="p-6 rounded-2xl bg-indigo-600/10 border border-indigo-500/20 text-center">
                 <p className="text-[10px] font-black text-indigo-400 uppercase tracking-[0.2em] mb-2">GPA ESTIMATION</p>
                 <p className="text-3xl font-black text-white leading-none">A+</p>
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
        <form onSubmit={handleSubmit} className="p-2 space-y-6">
          <div className="space-y-2">
            <label className="block text-[11px] font-black text-slate-500 uppercase tracking-[0.15em] leading-none mb-1">Task / Assignment</label>
            <select
              value={formData.task_id}
              onChange={e => setFormData({...formData, task_id: e.target.value})}
              className="input-dark leading-relaxed"
              required
            >
              <option value="" disabled className="bg-slate-900">Select Task</option>
              {tasks.map(t => (
                <option key={t.Id || t.id} value={t.Id || t.id} className="bg-slate-900 text-white">{t.Name_Task || t.name_task || 'Unknown Task'}</option>
              ))}
            </select>
          </div>
          
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
            <div className="space-y-2">
              <Input 
                label="Tanggal" 
                type="date"
                value={formData.tanggal} 
                onChange={e => setFormData({...formData, tanggal: e.target.value})} 
                required 
              />
            </div>
            <div className="space-y-2">
              <Input 
                label="Grades (0-100)" 
                type="number"
                min="0"
                max="100"
                value={formData.grades} 
                onChange={e => setFormData({...formData, grades: e.target.value})} 
                required 
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="block text-[11px] font-black text-slate-500 uppercase tracking-[0.15em] leading-none mb-1">Feedback Information</label>
            <textarea
              value={formData.keterangan}
              onChange={e => setFormData({...formData, keterangan: e.target.value})}
              required
              rows={4}
              placeholder="Berikan feedback akademik..."
              className="input-dark resize-none leading-relaxed py-4"
            />
          </div>
          
          <div className="pt-6 flex justify-end gap-4 border-t border-white/5">
            <button type="button" onClick={() => setIsModalOpen(false)}
              className="px-6 py-3 text-sm font-bold rounded-xl bg-white/5 hover:bg-white/10 text-slate-300 transition-colors leading-none">
              Cancel
            </button>
            <button type="submit" disabled={formLoading}
              className="btn-primary shadow-lg shadow-indigo-500/20"
              style={{ opacity: formLoading ? 0.6 : 1 }}>
              {formLoading ? 'Saving...' : 'Save Grade'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
