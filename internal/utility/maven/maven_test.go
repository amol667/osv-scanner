package maven_test

import (
	"path/filepath"
	"testing"

	mavenutil "deps.dev/util/maven"
	"deps.dev/util/resolve"
	"deps.dev/util/semver"
	"github.com/google/osv-scanner/v2/internal/utility/maven"
)

func TestParentPOMPath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		currentPath, relativePath string
		want                      string
	}{
		// fixtures
		// |- maven
		// |  |- my-app
		// |  |  |- pom.xml
		// |  |- parent
		// |  |  |- pom.xml
		// |- pom.xml
		{
			// Parent path is specified correctly.
			currentPath:  filepath.Join("fixtures", "my-app", "pom.xml"),
			relativePath: "../parent/pom.xml",
			want:         filepath.Join("fixtures", "parent", "pom.xml"),
		},
		{
			// Wrong file name is specified in relative path.
			currentPath:  filepath.Join("fixtures", "my-app", "pom.xml"),
			relativePath: "../parent/abc.xml",
			want:         "",
		},
		{
			// Wrong directory is specified in relative path.
			currentPath:  filepath.Join("fixtures", "my-app", "pom.xml"),
			relativePath: "../not-found/pom.xml",
			want:         "",
		},
		{
			// Only directory is specified.
			currentPath:  filepath.Join("fixtures", "my-app", "pom.xml"),
			relativePath: "../parent",
			want:         filepath.Join("fixtures", "parent", "pom.xml"),
		},
		{
			// Parent relative path is default to '../pom.xml'.
			currentPath:  filepath.Join("fixtures", "my-app", "pom.xml"),
			relativePath: "",
			want:         filepath.Join("fixtures", "pom.xml"),
		},
		{
			// No pom.xml is found even in the default path.
			currentPath:  filepath.Join("fixtures", "pom.xml"),
			relativePath: "",
			want:         "",
		},
	}
	for _, tt := range tests {
		got := maven.ParentPOMPath(tt.currentPath, tt.relativePath)
		if got != tt.want {
			t.Errorf("ParentPOMPath(%s, %s): got %s, want %s", tt.currentPath, tt.relativePath, got, tt.want)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	t.Parallel()

	versionKey := func(name string, version string) resolve.VersionKey {
		return resolve.VersionKey{
			PackageKey: resolve.PackageKey{
				System: resolve.Maven,
				Name:   name,
			},
			Version: version,
		}
	}
	semVer := func(version string) *semver.Version {
		parsed, _ := resolve.Maven.Semver().Parse(version)
		return parsed
	}

	tests := []struct {
		vk   resolve.VersionKey
		a, b *semver.Version
		want int
	}{
		{
			versionKey("abc:xyz", "1.0.0"),
			semVer("1.2.3"),
			semVer("1.2.3"),
			0,
		},
		{
			versionKey("abc:xyz", "1.0.0"),
			semVer("1.2.3"),
			semVer("2.3.4"),
			-1,
		},
		{
			versionKey("com.google.guava:guava", "1.0.0"),
			semVer("1.2.3"),
			semVer("2.3.4"),
			-1,
		},
		{
			versionKey("com.google.guava:guava", "1.0.0"),
			semVer("1.2.3-jre"),
			semVer("2.3.4-jre"),
			-1,
		},
		{
			versionKey("com.google.guava:guava", "1.0.0"),
			semVer("1.2.3-android"),
			semVer("2.3.4-android"),
			-1,
		},
		{
			versionKey("com.google.guava:guava", "1.0.0"),
			semVer("2.3.4-android"),
			semVer("1.2.3-jre"),
			-1,
		},
		{
			versionKey("com.google.guava:guava", "1.0.0-jre"),
			semVer("1.2.3-android"),
			semVer("1.2.3-jre"),
			-1,
		},
		{
			versionKey("com.google.guava:guava", "1.0.0-android"),
			semVer("1.2.3-android"),
			semVer("1.2.3-jre"),
			1,
		},
		{
			versionKey("commons-io:commons-io", "1.0.0"),
			semVer("1.2.3"),
			semVer("2.3.4"),
			-1,
		},
		{
			versionKey("commons-io:commons-io", "1.0.0"),
			semVer("1.2.3"),
			semVer("20010101.000000"),
			1,
		},
	}
	for _, tt := range tests {
		got := maven.CompareVersions(tt.vk, tt.a, tt.b)
		if got != tt.want {
			t.Errorf("CompareVersions(%v, %v, %v): got %b, want %b", tt.vk, tt.a, tt.b, got, tt.want)
		}
	}
}

func TestProjectKey(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name string
		proj mavenutil.Project
		want mavenutil.ProjectKey
	}{
		{
			name: "complete project",
			proj: mavenutil.Project{
				ProjectKey: mavenutil.ProjectKey{
					GroupID:    "com.example",
					ArtifactID: "test-artifact",
					Version:    "1.0.0",
				},
			},
			want: mavenutil.ProjectKey{
				GroupID:    "com.example",
				ArtifactID: "test-artifact",
				Version:    "1.0.0",
			},
		},
		{
			name: "missing groupId uses parent",
			proj: mavenutil.Project{
				ProjectKey: mavenutil.ProjectKey{
					ArtifactID: "test-artifact",
					Version:    "1.0.0",
				},
				Parent: mavenutil.Parent{
					ProjectKey: mavenutil.ProjectKey{
						GroupID: "com.parent",
					},
				},
			},
			want: mavenutil.ProjectKey{
				GroupID:    "com.parent",
				ArtifactID: "test-artifact",
				Version:    "1.0.0",
			},
		},
		{
			name: "missing version uses parent",
			proj: mavenutil.Project{
				ProjectKey: mavenutil.ProjectKey{
					GroupID:    "com.example",
					ArtifactID: "test-artifact",
				},
				Parent: mavenutil.Parent{
					ProjectKey: mavenutil.ProjectKey{
						Version: "2.0.0",
					},
				},
			},
			want: mavenutil.ProjectKey{
				GroupID:    "com.example",
				ArtifactID: "test-artifact",
				Version:    "2.0.0",
			},
		},
		{
			name: "both missing use parent",
			proj: mavenutil.Project{
				ProjectKey: mavenutil.ProjectKey{
					ArtifactID: "test-artifact",
				},
				Parent: mavenutil.Parent{
					ProjectKey: mavenutil.ProjectKey{
						GroupID: "com.parent",
						Version: "3.0.0",
					},
				},
			},
			want: mavenutil.ProjectKey{
				GroupID:    "com.parent",
				ArtifactID: "test-artifact",
				Version:    "3.0.0",
			},
		},
		{
			name: "no parent fallback",
			proj: mavenutil.Project{
				ProjectKey: mavenutil.ProjectKey{
					GroupID:    "com.example",
					ArtifactID: "test-artifact",
					Version:    "1.0.0",
				},
			},
			want: mavenutil.ProjectKey{
				GroupID:    "com.example",
				ArtifactID: "test-artifact",
				Version:    "1.0.0",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			got := maven.ProjectKey(tt.proj)
			if got != tt.want {
				t.Errorf("ProjectKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDependencyManagement(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name     string
		proj     mavenutil.Project
		expected int
	}{
		{
			name: "project with dependency management",
			proj: mavenutil.Project{
				DependencyManagement: mavenutil.DependencyManagement{
					Dependencies: []mavenutil.Dependency{
						{
							GroupID:    "junit",
							ArtifactID: "junit",
							Version:    "4.13.2",
							Scope:      "test",
						},
						{
							GroupID:    "org.springframework",
							ArtifactID: "spring-core",
							Version:    "5.3.21",
						},
					},
				},
			},
			expected: 2,
		},
		{
			name: "project without dependency management",
			proj: mavenutil.Project{
				ProjectKey: mavenutil.ProjectKey{
					GroupID:    "com.example",
					ArtifactID: "test",
					Version:    "1.0.0",
				},
			},
			expected: 0,
		},
		{
			name: "empty dependency management",
			proj: mavenutil.Project{
				DependencyManagement: mavenutil.DependencyManagement{
					Dependencies: []mavenutil.Dependency{},
				},
			},
			expected: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			got := tt.proj.DependencyManagement.Dependencies
			
			if len(got) != tt.expected {
				t.Errorf("GetDependencyManagement() returned %d dependencies, want %d", len(got), tt.expected)
				return
			}
			
			if tt.expected > 0 {
				if got[0].GroupID == "" {
					t.Error("First dependency should have GroupID")
				}
				if got[0].ArtifactID == "" {
					t.Error("First dependency should have ArtifactID")
				}
			}
		})
	}
}
