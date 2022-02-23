package main

func GetIAMPolicy() IAMPolicy {
	return IAMPolicy{
		Version: "something",
		Statement: []PolicyStatement{
			{
				Effect:   "Allow",
				Action:   []string{"ec2:CreateTags", "ec2:DeleteTags"},
				Resource: "arn:aws:ec2:*:*:security-group/*",
				Condition: map[string]map[string]string{
					"Null": {
						"aws:RequestTag/elbv2.k8s.aws/cluster":  "true",
						"aws:ResourceTag/elbv2.k8s.aws/cluster": "false",
					},
				},
			},
			{
				Effect:   "Allow",
				Action:   []string{"ec2:CreateTags", "ec2:DeleteTags"},
				Resource: "arn:aws:ec2:*:*:security-group/*",
				Condition: map[string]map[string]string{
					"Null": {
						"aws:RequestTag/elbv2.k8s.aws/cluster":  "true",
						"aws:ResourceTag/elbv2.k8s.aws/cluster": "false",
					},
				},
			},
		},
	}
}
