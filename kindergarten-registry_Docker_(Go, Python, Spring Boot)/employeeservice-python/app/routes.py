from flask import Blueprint, request, jsonify
import traceback
from .database import MongoDB
from .models import Employee

employee_bp = Blueprint('employee', __name__)

@employee_bp.route('/add-employee', methods=['POST', 'OPTIONS'])
def add_employee():
    if request.method == 'OPTIONS':
        return '', 200
    
    try:
        data = request.get_json()
        
        # Validate required fields
        if not all(k in data for k in ['name', 'id', 'position']):
            return jsonify({"error": "Missing required fields"}), 400
        
        employee = Employee.from_dict(data)
        
        # Check if employee already exists
        collection = MongoDB.get_collection('employees')
        existing = collection.find_one({"id": employee.id})
        if existing:
            return jsonify({"error": "Employee with this ID already exists"}), 409
        
        # Insert new employee
        result = collection.insert_one(employee.to_dict())
        
        return jsonify({
            "message": "Employee added successfully",
            "employee": employee.to_dict(),
            "inserted_id": str(result.inserted_id)
        }), 201
        
    except Exception as e:
        print(f"Error adding employee: {str(e)}")
        print(traceback.format_exc())
        return jsonify({"error": str(e)}), 500

@employee_bp.route('/employees', methods=['GET', 'OPTIONS'])
def get_employees():
    if request.method == 'OPTIONS':
        return '', 200
    
    try:
        collection = MongoDB.get_collection('employees')
        employees = list(collection.find({}, {'_id': 0}))
        return jsonify(employees), 200
    except Exception as e:
        print(f"Error getting employees: {str(e)}")
        return jsonify({"error": str(e)}), 500

@employee_bp.route('/delete-employee', methods=['DELETE', 'OPTIONS'])
def delete_employee():
    if request.method == 'OPTIONS':
        return '', 200
    
    try:
        employee_id = request.args.get('id')
        if not employee_id:
            return jsonify({"error": "ID parameter missing"}), 400
        
        collection = MongoDB.get_collection('employees')
        result = collection.delete_one({"id": employee_id})
        
        if result.deleted_count == 0:
            return jsonify({"error": "Employee not found"}), 404
        
        return jsonify({"message": "Employee deleted successfully"}), 200
    except Exception as e:
        print(f"Error deleting employee: {str(e)}")
        return jsonify({"error": str(e)}), 500

@employee_bp.route('/update-employee', methods=['PUT', 'OPTIONS'])
def update_employee():
    if request.method == 'OPTIONS':
        return '', 200
    
    try:
        data = request.get_json()
        
        if 'id' not in data:
            return jsonify({"error": "ID is required"}), 400
        
        collection = MongoDB.get_collection('employees')
        
        # Check if employee exists
        existing = collection.find_one({"id": data['id']})
        if not existing:
            return jsonify({"error": "Employee not found"}), 404
        
        # Update employee
        update_data = {k: v for k, v in data.items() if k in ['name', 'position']}
        result = collection.update_one(
            {"id": data['id']},
            {"$set": update_data}
        )
        
        if result.modified_count == 0:
            return jsonify({"message": "No changes made"}), 200
        
        # Return updated employee
        updated = collection.find_one({"id": data['id']}, {'_id': 0})
        return jsonify(updated), 200
        
    except Exception as e:
        print(f"Error updating employee: {str(e)}")
        return jsonify({"error": str(e)}), 500