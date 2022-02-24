package awsloadbalancercontroller

import cco "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"

type IAMPolicy struct {
    Version   string
    Statement []cco.StatementEntry
}

func GetIAMPolicy() IAMPolicy {
    return IAMPolicy{Version: "2012-10-17", Statement: []cco.StatementEntry{{Effect: "Allow", Action: []string{"elasticloadbalancing:AddTags", "elasticloadbalancing:RemoveTags"}, Resource: "arn:aws:elasticloadbalancing:*:*:targetgroup/*/*", PolicyCondition: cco.IAMPolicyCondition{"Null": cco.IAMPolicyConditionKeyValue{"aws:ResourceTag/elbv2.k8s.aws/cluster": "false", "aws:RequestTag/elbv2.k8s.aws/cluster": "true"}}}, {Effect: "Allow", Action: []string{"elasticloadbalancing:AddTags", "elasticloadbalancing:RemoveTags"}, Resource: "arn:aws:elasticloadbalancing:*:*:loadbalancer/net/*/*", PolicyCondition: cco.IAMPolicyCondition{"Null": cco.IAMPolicyConditionKeyValue{"aws:RequestTag/elbv2.k8s.aws/cluster": "true", "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"}}}, {Effect: "Allow", Action: []string{"elasticloadbalancing:AddTags", "elasticloadbalancing:RemoveTags"}, Resource: "arn:aws:elasticloadbalancing:*:*:loadbalancer/app/*/*", PolicyCondition: cco.IAMPolicyCondition{"Null": cco.IAMPolicyConditionKeyValue{"aws:RequestTag/elbv2.k8s.aws/cluster": "true", "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"}}}}}
}
