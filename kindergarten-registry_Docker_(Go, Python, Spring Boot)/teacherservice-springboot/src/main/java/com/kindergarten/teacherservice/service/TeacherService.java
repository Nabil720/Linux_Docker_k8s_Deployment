package com.kindergarten.teacherservice.service;

import com.kindergarten.teacherservice.model.Teacher;
import com.kindergarten.teacherservice.repository.TeacherRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import java.util.List;
import java.util.Optional;

@Service
public class TeacherService {
    
    @Autowired
    private TeacherRepository teacherRepository;
    
    public List<Teacher> getAllTeachers() {
        return teacherRepository.findAll();
    }
    
    public Teacher addTeacher(Teacher teacher) {
        // Check if teacher with same ID already exists
        if (teacherRepository.existsByTeacherId(teacher.getTeacherId())) {
            throw new RuntimeException("Teacher with this ID already exists");
        }
        return teacherRepository.save(teacher);
    }
    
    public Teacher updateTeacher(String teacherId, Teacher updatedTeacher) {
        Optional<Teacher> existingTeacher = teacherRepository.findByTeacherId(teacherId);
        
        if (existingTeacher.isPresent()) {
            Teacher teacher = existingTeacher.get();
            teacher.setName(updatedTeacher.getName());
            teacher.setSubject(updatedTeacher.getSubject());
            return teacherRepository.save(teacher);
        }
        throw new RuntimeException("Teacher not found");
    }
    
    public void deleteTeacher(String teacherId) {
        if (!teacherRepository.existsByTeacherId(teacherId)) {
            throw new RuntimeException("Teacher not found");
        }
        teacherRepository.deleteByTeacherId(teacherId);
    }
    
    public Optional<Teacher> getTeacherById(String teacherId) {
        return teacherRepository.findByTeacherId(teacherId);
    }
}