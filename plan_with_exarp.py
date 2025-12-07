#!/usr/bin/env python3
"""
Use exarp tools to generate planning for devwisdom-go
"""

import sys
from pathlib import Path

# Add parent project to path
sys.path.insert(0, str(Path(__file__).parent.parent))

from project_management_automation.tools.prd_generator import generate_prd
from project_management_automation.tools.task_discovery import discover_tasks_from_markdown
from project_management_automation.tools.task_analysis import analyze_task_structure

def main():
    devwisdom_path = Path(__file__).parent
    goals_path = devwisdom_path / "PROJECT_GOALS.md"
    todo_path = devwisdom_path / "TODO.md"
    output_path = devwisdom_path / "PLAN_OUTPUT.md"
    
    print("=" * 60)
    print("devwisdom-go Planning with Exarp Tools")
    print("=" * 60)
    print()
    
    # Read PROJECT_GOALS.md
    if goals_path.exists():
        print(f"📋 Reading PROJECT_GOALS.md...")
        with open(goals_path) as f:
            goals_content = f.read()
        
        # Generate PRD
        print(f"📝 Generating PRD...")
        try:
            prd_result = generate_prd(
                project_root=str(devwisdom_path.parent),
                goals_content=goals_content
            )
            print(f"✅ PRD generated")
        except Exception as e:
            print(f"⚠️  PRD generation failed: {e}")
            prd_result = None
    
    # Discover tasks from TODO.md
    if todo_path.exists():
        print(f"📋 Discovering tasks from TODO.md...")
        try:
            tasks = discover_tasks_from_markdown(str(todo_path))
            print(f"✅ Discovered {len(tasks)} tasks")
        except Exception as e:
            print(f"⚠️  Task discovery failed: {e}")
            tasks = []
    
    # Analyze task structure
    print(f"📊 Analyzing task structure...")
    try:
        analysis = analyze_task_structure()
        print(f"✅ Task analysis complete")
    except Exception as e:
        print(f"⚠️  Task analysis failed: {e}")
        analysis = None
    
    # Write output
    with open(output_path, 'w') as f:
        f.write("# devwisdom-go Planning Output\n\n")
        f.write("Generated using exarp tools\n\n")
        if prd_result:
            f.write("## PRD\n\n")
            f.write(str(prd_result))
            f.write("\n\n")
        if tasks:
            f.write(f"## Discovered Tasks ({len(tasks)})\n\n")
            for task in tasks:
                f.write(f"- {task}\n")
            f.write("\n")
    
    print()
    print(f"✅ Planning output written to: {output_path}")
    print()
    print("📌 Next Steps:")
    print("1. Review PLAN_OUTPUT.md")
    print("2. Use analyze_todo2_alignment to check alignment")
    print("3. Create structured tasks using agentic-tools MCP")

if __name__ == "__main__":
    main()
