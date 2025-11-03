package svnbranch

import "testing"

func TestGetSvnBranchName(t *testing.T) {
	type args struct {
		filePath string
		rootPath []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				filePath: "branches/feature/2021-05-01/bugfix/2021-05-01",
				rootPath: []string{"2021-05-01"},
			},
			want: "feature",
		},
		{
			name: "test2",
			args: args{
				filePath: "branches/feature/2021-05-01/bugfix/2021-05-01",
				rootPath: []string{"feature/2021-05-01/"},
			},
			want: "",
		},
		{
			name: "test4",
			args: args{
				filePath: "trunk/feature/2021-05-01/bugfix/2021-05-01",
				rootPath: []string{"feature/2021-05-01/"},
			},
			want: "trunk",
		},
		{
			name: "test5",
			args: args{
				filePath: "trunk/feature/2021-05-01/bugfix/2021-05-01",
			},
			want: "trunk",
		},
		{
			name: "test6",
			args: args{
				filePath: "branches/feature/2021-05-01/bugfix/2021-05-01",
				rootPath: []string{"/2021-05-01/"},
			},
			want: "feature",
		},
		{
			name: "test7",
			args: args{
				filePath: "branches/feature/2021-05-01/bugfix/2021-05-01",
			},
			want: "feature",
		},
		{
			name: "test8",
			args: args{
				filePath: "branches/feature/feature_dev/bugfix/2021-05-01",
			},
			want: "feature_dev",
		},
		{
			name: "test9",
			args: args{
				filePath: "branches/",
			},
			want: "",
		},
		{
			name: "test10",
			args: args{
				filePath: "branches/obt_v2/obt_v2_oversea_online/",
			},
			want: "obt_v2_oversea_online",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := GetSvnBranchNameByPath(tt.args.filePath, tt.args.rootPath...); got != tt.want {
					t.Errorf("GetSvnBranchNameByPath() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
