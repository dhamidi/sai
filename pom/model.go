package pom

import "encoding/xml"

type Project struct {
	XMLName                xml.Name                `xml:"project"`
	ModelVersion           string                  `xml:"modelVersion"`
	GroupID                string                  `xml:"groupId"`
	ArtifactID             string                  `xml:"artifactId"`
	Version                string                  `xml:"version"`
	Packaging              string                  `xml:"packaging"`
	Name                   string                  `xml:"name"`
	Description            string                  `xml:"description"`
	URL                    string                  `xml:"url"`
	InceptionYear          string                  `xml:"inceptionYear"`
	Parent                 *Parent                 `xml:"parent"`
	Modules                []string                `xml:"modules>module"`
	Properties             *Properties             `xml:"properties"`
	Dependencies           []Dependency            `xml:"dependencies>dependency"`
	DependencyManagement   *DependencyManagement   `xml:"dependencyManagement"`
	Build                  *Build                  `xml:"build"`
	Reporting              *Reporting              `xml:"reporting"`
	Licenses               []License               `xml:"licenses>license"`
	Organization           *Organization           `xml:"organization"`
	Developers             []Developer             `xml:"developers>developer"`
	Contributors           []Contributor           `xml:"contributors>contributor"`
	IssueManagement        *IssueManagement        `xml:"issueManagement"`
	CIManagement           *CIManagement           `xml:"ciManagement"`
	MailingLists           []MailingList           `xml:"mailingLists>mailingList"`
	SCM                    *SCM                    `xml:"scm"`
	Prerequisites          *Prerequisites          `xml:"prerequisites"`
	Repositories           []Repository            `xml:"repositories>repository"`
	PluginRepositories     []Repository            `xml:"pluginRepositories>pluginRepository"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement"`
	Profiles               []Profile               `xml:"profiles>profile"`
}

type Parent struct {
	GroupID      string `xml:"groupId"`
	ArtifactID   string `xml:"artifactId"`
	Version      string `xml:"version"`
	RelativePath string `xml:"relativePath"`
}

type Properties struct {
	Entries map[string]string
}

func (p *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	p.Entries = make(map[string]string)
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			var value string
			if err := d.DecodeElement(&value, &t); err != nil {
				return err
			}
			p.Entries[t.Name.Local] = value
		case xml.EndElement:
			if t.Name == start.Name {
				return nil
			}
		}
	}
}

type Dependency struct {
	GroupID    string      `xml:"groupId"`
	ArtifactID string      `xml:"artifactId"`
	Version    string      `xml:"version"`
	Type       string      `xml:"type"`
	Classifier string      `xml:"classifier"`
	Scope      string      `xml:"scope"`
	SystemPath string      `xml:"systemPath"`
	Optional   string      `xml:"optional"`
	Exclusions []Exclusion `xml:"exclusions>exclusion"`
}

type Exclusion struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
}

type DependencyManagement struct {
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

type Build struct {
	SourceDirectory       string            `xml:"sourceDirectory"`
	ScriptSourceDirectory string            `xml:"scriptSourceDirectory"`
	TestSourceDirectory   string            `xml:"testSourceDirectory"`
	OutputDirectory       string            `xml:"outputDirectory"`
	TestOutputDirectory   string            `xml:"testOutputDirectory"`
	DefaultGoal           string            `xml:"defaultGoal"`
	Directory             string            `xml:"directory"`
	FinalName             string            `xml:"finalName"`
	Filters               []string          `xml:"filters>filter"`
	Resources             []Resource        `xml:"resources>resource"`
	TestResources         []Resource        `xml:"testResources>testResource"`
	Plugins               []Plugin          `xml:"plugins>plugin"`
	PluginManagement      *PluginManagement `xml:"pluginManagement"`
	Extensions            []Extension       `xml:"extensions>extension"`
}

type Resource struct {
	TargetPath string   `xml:"targetPath"`
	Filtering  string   `xml:"filtering"`
	Directory  string   `xml:"directory"`
	Includes   []string `xml:"includes>include"`
	Excludes   []string `xml:"excludes>exclude"`
}

type Plugin struct {
	GroupID       string            `xml:"groupId"`
	ArtifactID    string            `xml:"artifactId"`
	Version       string            `xml:"version"`
	Extensions    string            `xml:"extensions"`
	Inherited     string            `xml:"inherited"`
	Configuration *Configuration    `xml:"configuration"`
	Dependencies  []Dependency      `xml:"dependencies>dependency"`
	Executions    []PluginExecution `xml:"executions>execution"`
}

type PluginExecution struct {
	ID            string         `xml:"id"`
	Phase         string         `xml:"phase"`
	Goals         []string       `xml:"goals>goal"`
	Inherited     string         `xml:"inherited"`
	Configuration *Configuration `xml:"configuration"`
}

type Configuration struct {
	Raw []byte `xml:",innerxml"`
}

type PluginManagement struct {
	Plugins []Plugin `xml:"plugins>plugin"`
}

type Extension struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type Reporting struct {
	ExcludeDefaults string         `xml:"excludeDefaults"`
	OutputDirectory string         `xml:"outputDirectory"`
	Plugins         []ReportPlugin `xml:"plugins>plugin"`
}

type ReportPlugin struct {
	GroupID       string         `xml:"groupId"`
	ArtifactID    string         `xml:"artifactId"`
	Version       string         `xml:"version"`
	Inherited     string         `xml:"inherited"`
	Configuration *Configuration `xml:"configuration"`
	ReportSets    []ReportSet    `xml:"reportSets>reportSet"`
}

type ReportSet struct {
	ID            string         `xml:"id"`
	Reports       []string       `xml:"reports>report"`
	Inherited     string         `xml:"inherited"`
	Configuration *Configuration `xml:"configuration"`
}

type License struct {
	Name         string `xml:"name"`
	URL          string `xml:"url"`
	Distribution string `xml:"distribution"`
	Comments     string `xml:"comments"`
}

type Organization struct {
	Name string `xml:"name"`
	URL  string `xml:"url"`
}

type Developer struct {
	ID              string      `xml:"id"`
	Name            string      `xml:"name"`
	Email           string      `xml:"email"`
	URL             string      `xml:"url"`
	Organization    string      `xml:"organization"`
	OrganizationURL string      `xml:"organizationUrl"`
	Roles           []string    `xml:"roles>role"`
	Timezone        string      `xml:"timezone"`
	Properties      *Properties `xml:"properties"`
}

type Contributor struct {
	Name            string      `xml:"name"`
	Email           string      `xml:"email"`
	URL             string      `xml:"url"`
	Organization    string      `xml:"organization"`
	OrganizationURL string      `xml:"organizationUrl"`
	Roles           []string    `xml:"roles>role"`
	Timezone        string      `xml:"timezone"`
	Properties      *Properties `xml:"properties"`
}

type IssueManagement struct {
	System string `xml:"system"`
	URL    string `xml:"url"`
}

type CIManagement struct {
	System    string     `xml:"system"`
	URL       string     `xml:"url"`
	Notifiers []Notifier `xml:"notifiers>notifier"`
}

type Notifier struct {
	Type          string         `xml:"type"`
	SendOnError   string         `xml:"sendOnError"`
	SendOnFailure string         `xml:"sendOnFailure"`
	SendOnSuccess string         `xml:"sendOnSuccess"`
	SendOnWarning string         `xml:"sendOnWarning"`
	Configuration *Configuration `xml:"configuration"`
}

type MailingList struct {
	Name          string   `xml:"name"`
	Subscribe     string   `xml:"subscribe"`
	Unsubscribe   string   `xml:"unsubscribe"`
	Post          string   `xml:"post"`
	Archive       string   `xml:"archive"`
	OtherArchives []string `xml:"otherArchives>otherArchive"`
}

type SCM struct {
	Connection          string `xml:"connection"`
	DeveloperConnection string `xml:"developerConnection"`
	Tag                 string `xml:"tag"`
	URL                 string `xml:"url"`
}

type Prerequisites struct {
	Maven string `xml:"maven"`
}

type Repository struct {
	ID        string            `xml:"id"`
	Name      string            `xml:"name"`
	URL       string            `xml:"url"`
	Layout    string            `xml:"layout"`
	Releases  *RepositoryPolicy `xml:"releases"`
	Snapshots *RepositoryPolicy `xml:"snapshots"`
}

type RepositoryPolicy struct {
	Enabled        string `xml:"enabled"`
	UpdatePolicy   string `xml:"updatePolicy"`
	ChecksumPolicy string `xml:"checksumPolicy"`
}

type DistributionManagement struct {
	Repository         *DeploymentRepository `xml:"repository"`
	SnapshotRepository *DeploymentRepository `xml:"snapshotRepository"`
	Site               *Site                 `xml:"site"`
	Relocation         *Relocation           `xml:"relocation"`
	DownloadURL        string                `xml:"downloadUrl"`
	Status             string                `xml:"status"`
}

type DeploymentRepository struct {
	UniqueVersion string `xml:"uniqueVersion"`
	ID            string `xml:"id"`
	Name          string `xml:"name"`
	URL           string `xml:"url"`
	Layout        string `xml:"layout"`
}

type Site struct {
	ID   string `xml:"id"`
	Name string `xml:"name"`
	URL  string `xml:"url"`
}

type Relocation struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
	Message    string `xml:"message"`
}

type Profile struct {
	ID                     string                  `xml:"id"`
	Activation             *Activation             `xml:"activation"`
	Build                  *Build                  `xml:"build"`
	Modules                []string                `xml:"modules>module"`
	Repositories           []Repository            `xml:"repositories>repository"`
	PluginRepositories     []Repository            `xml:"pluginRepositories>pluginRepository"`
	Dependencies           []Dependency            `xml:"dependencies>dependency"`
	Reporting              *Reporting              `xml:"reporting"`
	DependencyManagement   *DependencyManagement   `xml:"dependencyManagement"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement"`
	Properties             *Properties             `xml:"properties"`
}

type Activation struct {
	ActiveByDefault string              `xml:"activeByDefault"`
	JDK             string              `xml:"jdk"`
	OS              *ActivationOS       `xml:"os"`
	Property        *ActivationProperty `xml:"property"`
	File            *ActivationFile     `xml:"file"`
}

type ActivationOS struct {
	Name    string `xml:"name"`
	Family  string `xml:"family"`
	Arch    string `xml:"arch"`
	Version string `xml:"version"`
}

type ActivationProperty struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type ActivationFile struct {
	Exists  string `xml:"exists"`
	Missing string `xml:"missing"`
}
